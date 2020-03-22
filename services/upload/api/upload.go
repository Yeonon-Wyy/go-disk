package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-disk/common"
	"go-disk/common/constant"
	"go-disk/common/log4disk"
	"go-disk/common/mqproto"
	"go-disk/common/utils"
	"go-disk/services/upload/config"
	"go-disk/services/upload/dao"
	"go-disk/services/upload/db"
	"go-disk/services/upload/midware/mq"
	redisconn "go-disk/services/upload/midware/redis"
	"go-disk/services/upload/vo"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	businessConfig = config.Conf.Business
	storeConfig    = config.Conf.Store
	mqConfig       = config.Conf.Mq
)

func UploadFile() gin.HandlerFunc {
	return func(context *gin.Context) {
		fh, err := context.FormFile("file")
		filePath := context.PostForm("file_path")
		username := context.PostForm("username")

		if filePath == "" || username == "" {
			log4disk.E("file path or username can't empty")
			context.JSON(http.StatusBadRequest, common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		if err != nil {
			log4disk.E("read from file error : %v", err)
			context.JSON(http.StatusInternalServerError, common.NewServiceResp(common.RespCodeReadFileError, nil))
			return
		}

		file, err := fh.Open()
		if err != nil {
			log4disk.E("open file error : %v", err)
			context.JSON(http.StatusInternalServerError, common.NewServiceResp(common.RespCodeOpenFileError, nil))
			return
		}
		defer file.Close()

		fileHash := utils.MultipartFileSha1(&file)
		filename := fh.Filename
		fileSize := fh.Size

		exist := db.ExistFile(fileHash)

		//如果唯一file表中不存在重复的文件，才需要拷贝上传到存储器中去
		if !exist {
			newFile, err := os.Create(businessConfig.FileStorePath + fh.Filename)
			if err != nil {
				log4disk.E("create file error : %v", err)
				context.JSON(http.StatusInternalServerError, common.NewServiceResp(common.RespCodeCreateFileError, nil))
				return
			}

			file.Seek(0, 0)
			newFile.Seek(0, 0)
			_, err = io.Copy(newFile, file)

			if err != nil {
				log4disk.E("copy file error : %v", err)
				context.JSON(http.StatusInternalServerError, common.NewServiceResp(common.RespCodeCopyFileError, nil))
				return
			}

			_, fName := filepath.Split(newFile.Name())

			//手动关闭 file，防止异步进程无法删除文件
			newFile.Close()


			//set file meta
			tblFile := dao.TableFileDao{
				FileHash: fileHash,
				Filename: fName,
				FileSize: fileSize,
				FileAddr: businessConfig.FileStorePath + fName,
				Status:   constant.FileStatusAvailable,
			}

			//同步到ceph中
			err = transFileToCeph(fileHash, tblFile.FileAddr)
			if err != nil {
				log4disk.E("put data to ceph error : %v", err)
				context.JSON(http.StatusInternalServerError,
					common.NewServiceResp(common.RespCodePutDataToCephError, nil))
				return
			}

			db.OnFileUploadFinished(
				tblFile.FileHash,
				tblFile.Filename,
				tblFile.FileAddr,
				tblFile.FileSize,
				tblFile.Status,
			)

		}

		//如果file已经上传过了，那么久需要修改文件名，避免重名，同时状态需要和唯一文件表里的状态一样
		if exist {
			idx := strings.Index(filename, ".")
			if idx == -1 {
				filename += strconv.FormatInt(time.Now().Unix(), 10)
			} else {
				filename = fmt.Sprintf("%s-%d.%s", filename[0:idx], time.Now().Unix(), filename[idx+1:])
			}
		}

		//写入userfile 表里
		status := db.GetStatus(fileHash)
		ok := db.InsertUserFile(username, fileHash, filename, filePath, fileSize, status)
		if !ok {
			log4disk.E("upload file failed, file hash is : %s", fileHash)
			context.JSON(http.StatusInternalServerError,
				common.NewServiceResp(common.RespCodeUploadFileError, nil))
			return
		}
		log4disk.E("upload file success, file hash is : %s", fileHash)
		context.JSON(http.StatusOK, common.NewServiceResp(common.RespCodeSuccess, nil))
	}
}

func TryFastUpload() gin.HandlerFunc {
	return func(context *gin.Context) {
		var req vo.FastUploadReq
		if err := context.ShouldBind(&req); err != nil {
			log4disk.E("bind request parameters error %v", err)
			context.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		if exist := db.ExistFile(req.FileHash); !exist {
			log4disk.E("not exist this file %s, failed to fast upload", req.FileHash)
			context.JSON(http.StatusNoContent,
				common.NewServiceResp(common.RespCodeFastUploadFailed, nil))
			return
		}

		ok := db.InsertUserFile(req.Username, req.FileHash, req.Filename, req.FilePath, req.FileSize, constant.FileStatusAvailable)
		if !ok {
			context.JSON(http.StatusInternalServerError,
				common.NewServiceResp(common.RespCodeUploadFileError, nil))
			return
		}

		context.JSON(http.StatusNoContent,
			common.NewServiceResp(common.RespCodeSuccess, nil))
	}
}

func InitialMultipartUpload() gin.HandlerFunc {
	return func(context *gin.Context) {
		var req vo.MPUploadInitReq
		if err := context.ShouldBind(&req); err != nil {
			log4disk.E("bind request parameters error %v", err)
			context.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		redisClient, err := redisconn.FSConn()
		if err != nil {
			log4disk.E("failed to connect redis server : %v", err)
			context.JSON(http.StatusInternalServerError,
				common.NewServiceResp(common.RespCodeConnectFSRedisServerError, nil))
			return
		}
		defer redisClient.Close()

		mpUploadInfo := vo.MPUploadInfo{
			FileHash:   req.FileHash,
			FileSize:   req.FileSize,
			UploadId:   req.Username + fmt.Sprintf("%x", time.Now().UnixNano()),
			ChunkSize:  constant.FileMPUploadChunkSize,
			ChunkCount: int(math.Ceil(float64(req.FileSize) / constant.FileMPUploadChunkSize)),
		}
		redisClient.HSet("MP_"+mpUploadInfo.UploadId, "chunk_count", mpUploadInfo.ChunkCount)
		redisClient.HSet("MP_"+mpUploadInfo.UploadId, "file_hash", mpUploadInfo.FileHash)
		redisClient.HSet("MP_"+mpUploadInfo.UploadId, "file_size", mpUploadInfo.FileSize)

		context.JSON(http.StatusOK,
			common.NewServiceResp(common.RespCodeSuccess, mpUploadInfo))
	}
}

func UploadPart() gin.HandlerFunc {
	return func(context *gin.Context) {
		uploadId := context.Query("uploadid")
		index := context.Query("index")
		if len(uploadId) == 0 || len(index) == 0 {
			log4disk.E("bind request parameters error")
			context.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		redisClient, err := redisconn.FSConn()
		if err != nil {
			log4disk.E("failed to connect redis server : %v", err)
			context.JSON(http.StatusInternalServerError,
				common.NewServiceResp(common.RespCodeConnectFSRedisServerError, nil))
			return
		}
		defer redisClient.Close()

		fpath := businessConfig.FileStorePath + uploadId + "/" + index
		os.MkdirAll(path.Dir(fpath), 0777)
		fd, err := os.Create(fpath)
		if err != nil {
			log4disk.E("create file error : %v", err)
			context.JSON(http.StatusInternalServerError,
				common.NewServiceResp(common.RespCodeOpenFileError, nil))
			return
		}
		defer fd.Close()
		//buf := make([]byte, 1024 * 1024)
		data, err := ioutil.ReadAll(context.Request.Body)
		if err != nil {
			log4disk.E("read request data error : %v", err)
			context.JSON(http.StatusInternalServerError,
				common.NewServiceResp(common.RespCodeReadDataError, nil))
			return
		}
		n, err := fd.Write(data)
		if err != nil {
			log4disk.E("write data to file error : %v", err)
			context.JSON(http.StatusInternalServerError,
				common.NewServiceResp(common.RespCodeWriteFileError, nil))
			return
		}
		log4disk.E("write to file success size : n = %d byte", n)

		redisClient.HSet("MP_"+uploadId, "chkidx_"+index, 1)

		context.JSON(http.StatusOK,
			common.NewServiceResp(common.RespCodeSuccess, nil))
	}
}

func CompleteUpload() gin.HandlerFunc {
	return func(context *gin.Context) {
		var req vo.MPUploadCompleteReq
		if err := context.ShouldBind(&req); err != nil {
			log4disk.E("bind request parameters error %v", err)
			context.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		redisClient, err := redisconn.FSConn()
		if err != nil {
			log4disk.E("failed to connect redis server : %v", err)
			context.JSON(http.StatusInternalServerError,
				common.NewServiceResp(common.RespCodeConnectFSRedisServerError, nil))
			return
		}
		defer redisClient.Close()

		result, err := redisClient.HGetAll("MP_" + req.UploadId).Result()
		if err != nil {
			log4disk.E("failed to get value from redis server : %v", err)
			context.JSON(http.StatusInternalServerError,
				common.NewServiceResp(common.RespCodeCompleteUploadError, nil))
			return
		}

		totalCount, chunkCount := 0, 0
		for key, value := range result {
			if key == "chunk_count" {
				totalCount, _ = strconv.Atoi(value)
			} else if strings.HasPrefix(key, "chkidx_") && value == "1" {
				chunkCount++
			}
		}

		if totalCount != chunkCount {
			log4disk.E("invaild request")
			context.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		//开始做分块合并操作，开启一个协程去做
		tempFileDirPath := businessConfig.FileStorePath + req.UploadId + "/"
		dstFileName := businessConfig.FileStorePath + req.Filename

		go func(fileHash, srcFileDirPath, dstFileName string, count int) {
			err := mergeFilePart(tempFileDirPath, dstFileName, count)
			if err != nil {
				log4disk.E("merge file error : %v", err)
				return
			}
			//文件合并完成就可以删除了保存分块文件的文件夹了
			os.RemoveAll(srcFileDirPath)

			//同步到ceph中
			err = transFileToCeph(fileHash, dstFileName)
			if err != nil {
				log4disk.E("put data to ceph error : %v", err)
				return
			}
		}(req.FileHash, tempFileDirPath, dstFileName, chunkCount)

		//写入数据库
		db.OnFileUploadFinished(
			req.FileHash,
			req.Filename,
			"",
			req.FileSize,
			constant.FileStatusAvailable,
			)

		db.InsertUserFile(
			req.Username,
			req.FileHash,
			req.Filename,
			req.FilePath,
			req.FileSize,
			constant.FileStatusAvailable)

		//合并和同步到ceph的操作不会用到redis里的数据，所以直接删除就行了
		redisClient.Del("MP_" + req.UploadId)

		context.JSON(http.StatusOK,
			common.NewServiceResp(common.RespCodeSuccess, nil))

	}
}

func transFileToCeph(fileHash string, fileLocation string) error {
	msgData := mqproto.RabbitMessage{
		FileHash:     fileHash,
		SrcLocation:  fileLocation,
		DstLocation:  storeConfig.Ceph.FilePathPrefix + fileHash,
		DstStoreType: common.StoreCeph,
	}
	msgDataJson, err := json.Marshal(msgData)
	if err != nil {
		return err
	}
	suc := mq.RabbitPublish(
		mqConfig.Rabbit.ExchangeName,
		mqConfig.Rabbit.CephRoutingKey,
		msgDataJson)

	if !suc {
		return errors.New("push message to rabbit mq error")
	}
	return nil
}

func mergeFilePart(tempFileDirPath string, dstFileName string, chunkCount int) error {
	dstFd, err := os.OpenFile(dstFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}
	defer dstFd.Close()

	for i := 1; i <= chunkCount; i++ {
		fileName := tempFileDirPath + strconv.Itoa(i)
		srcFd, err := os.OpenFile(fileName, os.O_RDONLY, os.ModePerm)
		if err != nil {
			return err
		}

		b, err := ioutil.ReadAll(srcFd)
		_, err = dstFd.Write(b)
		if err != nil {
			return err
		}
		srcFd.Close()
	}

	return nil
}
