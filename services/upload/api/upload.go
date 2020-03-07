package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	redisconn "go-disk/cache/redis"
	"go-disk/common"
	"go-disk/common/constant"
	commonconfig "go-disk/config"
	"go-disk/mq"
	"go-disk/services/upload/config"
	uploadconfig "go-disk/services/upload/config"

	"go-disk/services/upload/dao"
	"go-disk/services/upload/db"
	"go-disk/services/upload/vo"
	"go-disk/utils"
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

func UploadFile() gin.HandlerFunc {
	return func(context *gin.Context) {
		fh, err := context.FormFile("file")

		if err != nil {
			log.Printf("read from file error : %v", err)
			context.JSON(http.StatusInternalServerError, common.NewServiceResp(common.RespCodeReadFileError, nil))
			return
		}

		file, err := fh.Open()
		if err != nil {
			log.Printf("open file error : %v", err)
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
			newFile, err := os.Create(uploadconfig.FileStoreDir + fh.Filename)
			if err != nil {
				log.Printf("create file error : %v", err)
				context.JSON(http.StatusInternalServerError, common.NewServiceResp(common.RespCodeCreateFileError, nil))
				return
			}

			file.Seek(0, 0)
			newFile.Seek(0, 0)
			_, err = io.Copy(newFile, file)

			if err != nil {
				log.Printf("copy file error : %v", err)
				context.JSON(http.StatusInternalServerError, common.NewServiceResp(common.RespCodeCopyFileError, nil))
				return
			}

			defer newFile.Close()

			_, fName := filepath.Split(newFile.Name())

			//set file meta
			tblFile := dao.TableFile{
				FileSha1:     fileHash,
				Filename:     fName,
				FileSize:     fileSize,
				FileLocation: config.FileStoreDir + fName,
				CreateAt:     time.Now(),
				UpdateAt:     time.Now(),
				Status:       constant.FileStatusAvailable,
			}

			//同步到ceph中
			err = transFileToCeph(fileHash, tblFile.FileLocation)
			if err != nil {
				log.Printf("put data to ceph error : %v", err)
				context.JSON(http.StatusInternalServerError,
					common.NewServiceResp(common.RespCodePutDataToCephError, nil))
				return
			}

			db.OnFileUploadFinished(
				tblFile.FileSha1,
				tblFile.Filename,
				tblFile.FileLocation,
				tblFile.FileSize,
				tblFile.Status,
				tblFile.UpdateAt,
				tblFile.UpdateAt,
				)

		}

		//如果file已经上传过了，那么久需要修改文件名，避免重名
		if exist {
			idx := strings.Index(filename, ".")
			if idx == -1 {
				filename += strconv.FormatInt(time.Now().Unix(), 10)
			} else {
				filename = fmt.Sprintf("%s-%d.%s", filename[0:idx], time.Now().Unix(), filename[idx+1:])
			}
		}

		//写入userfile 表里
		username := context.Query("username")
		ok := db.InsertUserFile(username, fileHash, filename, fileSize)
		if !ok {
			log.Printf("upload file failed, file hash is : %s", fileHash)
			context.JSON(http.StatusInternalServerError,
				common.NewServiceResp(common.RespCodeUploadFileError, nil))
			return
		}
		log.Printf("upload file success, file hash is : %s", fileHash)
		context.JSON(http.StatusOK, common.NewServiceResp(common.RespCodeSuccess, nil))
	}
}

func TryFastUpload() gin.HandlerFunc {
	return func(context *gin.Context) {
		var req vo.FastUploadReq
		if err := context.ShouldBind(&req); err != nil {
			log.Printf("bind request parameters error %v", err)
			context.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		if exist := db.ExistFile(req.FileHash); !exist {
			log.Printf("not exist this file %s, failed to fast upload", req.FileHash)
			context.JSON(http.StatusNoContent,
				common.NewServiceResp(common.RespCodeFastUploadFailed, nil))
			return
		}

		ok := db.InsertUserFile(req.Username, req.FileHash, req.Filename, req.FileSize)
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
			log.Printf("bind request parameters error %v", err)
			context.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		redisClient, err := redisconn.FSConn()
		if err != nil {
			log.Printf("failed to connect redis server : %v", err)
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
			ChunkCount: int(math.Ceil(float64(req.FileSize)/constant.FileMPUploadChunkSize)),
		}
		redisClient.HSet("MP_" + mpUploadInfo.UploadId, "chunk_count", mpUploadInfo.ChunkCount)
		redisClient.HSet("MP_" + mpUploadInfo.UploadId, "file_hash", mpUploadInfo.FileHash)
		redisClient.HSet("MP_" + mpUploadInfo.UploadId, "file_size", mpUploadInfo.FileSize)

		context.JSON(http.StatusOK,
			common.NewServiceResp(common.RespCodeSuccess, mpUploadInfo))
	}
}

func UploadPart() gin.HandlerFunc {
	return func(context *gin.Context) {
		uploadId := context.Query("uploadid")
		index := context.Query("index")
		if len(uploadId) == 0 || len(index) == 0 {
			log.Printf("bind request parameters error")
			context.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		redisClient, err := redisconn.FSConn()
		if err != nil {
			log.Printf("failed to connect redis server : %v", err)
			context.JSON(http.StatusInternalServerError,
				common.NewServiceResp(common.RespCodeConnectFSRedisServerError, nil))
			return
		}
		defer redisClient.Close()

		fpath := config.FileStoreDir + uploadId + "/" + index
		os.MkdirAll(path.Dir(fpath), 0777)
		fd, err := os.Create(fpath)
		if err != nil {
			log.Printf("create file error : %v", err)
			context.JSON(http.StatusInternalServerError,
				common.NewServiceResp(common.RespCodeOpenFileError, nil))
			return
		}
		defer fd.Close()
		//buf := make([]byte, 1024 * 1024)
		data, err := ioutil.ReadAll(context.Request.Body)
		if err != nil {
			log.Printf("read request data error : %v", err)
			context.JSON(http.StatusInternalServerError,
				common.NewServiceResp(common.RespCodeReadDataError, nil))
			return
		}
		n, err := fd.Write(data)
		if err != nil {
			log.Printf("write data to file error : %v", err)
			context.JSON(http.StatusInternalServerError,
				common.NewServiceResp(common.RespCodeWriteFileError, nil))
			return
		}
		log.Printf("write to file success size : n = %d byte", n)


		redisClient.HSet("MP_" + uploadId, "chkidx_" + index, 1)

		context.JSON(http.StatusOK,
			common.NewServiceResp(common.RespCodeSuccess, nil))
	}
}

func CompleteUpload() gin.HandlerFunc {
	return func(context *gin.Context) {
		var req vo.MPUploadCompleteReq
		if err := context.ShouldBind(&req); err != nil {
			log.Printf("bind request parameters error %v", err)
			context.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		redisClient, err := redisconn.FSConn()
		if err != nil {
			log.Printf("failed to connect redis server : %v", err)
			context.JSON(http.StatusInternalServerError,
				common.NewServiceResp(common.RespCodeConnectFSRedisServerError, nil))
			return
		}
		defer redisClient.Close()

		result, err := redisClient.HGetAll("MP_" + req.UploadId).Result()
		if err != nil {
			log.Printf("failed to get value from redis server : %v", err)
			context.JSON(http.StatusInternalServerError,
				common.NewServiceResp(common.RespCodeCompleteUploadError, nil))
			return
		}

		totalCount, chunkCount := 0, 0
		for key, value := range result {
			if key == "chunk_count" {
				totalCount,_ = strconv.Atoi(value)
			} else if strings.HasPrefix(key, "chkidx_") && value == "1" {
				chunkCount++
			}
		}

		if totalCount != chunkCount {
			log.Printf("invaild request")
			context.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		//TODO: 需要分块合并操作，文件名应该是config.FileStoreDir + req.Filename

		//同步到ceph中
		//TODO: 暂时确定文件名为这个，需要和上述合并操作联调
		err = transFileToCeph(req.FileHash, config.FileStoreDir + req.Filename)
		if err != nil {
			log.Printf("put data to ceph error : %v", err)
			context.JSON(http.StatusInternalServerError,
				common.NewServiceResp(common.RespCodePutDataToCephError, nil))
			return
		}

		//写入数据库
		db.OnFileUploadFinished(
			req.FileHash,
			req.Filename,
			"",
			req.FileSize,
			constant.UserStatusAvailable,
			time.Now(),
			time.Now())

		db.InsertUserFile(
			req.Username,
			req.FileHash,
			req.Filename,
			req.FileSize)

		redisClient.Del("MP_" + req.UploadId)

		//delete temp dir
		fpath := config.FileStoreDir + req.UploadId + "/"
		os.RemoveAll(fpath)

		context.JSON(http.StatusOK,
			common.NewServiceResp(common.RespCodeSuccess, nil))

	}
}

func transFileToCeph(fileHash string, fileLocation string) error {
	msgData := mq.RabbitMessage{
		FileHash:     fileHash,
		SrcLocation:  fileLocation,
		DstLocation:  commonconfig.CephFilePathPrefix + fileHash,
		DstStoreType: common.StoreCeph,
	}
	msgDataJson, err := json.Marshal(msgData)
	if err != nil {
		return err
	}
	suc := mq.RabbitPublish(
		commonconfig.RabbitExchangeName,
		commonconfig.RabbitCephRoutingKey,
		msgDataJson)

	if !suc {
		return errors.New("push message to rabbit mq error")
	}
	return nil
}
