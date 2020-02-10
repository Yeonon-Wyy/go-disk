package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-disk/auth"
	"go-disk/common"
	"go-disk/common/constant"
	"go-disk/config"
	"go-disk/db"
	"go-disk/meta"
	"go-disk/model"
	"go-disk/utils"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type FilesServiceHandler struct {
	BashPath string
}


func (f FilesServiceHandler) Init(group *gin.RouterGroup) {


	group.Use(auth.AuthorizeInterceptor())
	group.StaticFile("/upload", "./static/view/index.html")
	group.POST("/upload", uploadFile())
	group.POST("/fastupload", tryFastUpload())

	group.GET("/meta", getFileMeta())

	group.PUT("/meta", updateFileMeta())
	group.POST("/meta", queryFileList())

	group.DELETE("/delete", deleteFile())
	group.POST("/download", downloadHandler())
}

func uploadFile() gin.HandlerFunc {
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
			newFile, err := os.Create(config.FileStoreDir + fh.Filename)
			if err != nil {
				log.Printf("create file error : %v", err)
				context.JSON(http.StatusInternalServerError, common.NewServiceResp(common.RespCodeCreateFileError, nil))
				return
			}

			_, err = io.Copy(newFile, file)

			if err != nil {
				log.Printf("copy file error : %v", err)
				context.JSON(http.StatusInternalServerError, common.NewServiceResp(common.RespCodeCopyFileError, nil))
				return
			}

			defer newFile.Close()
			newFile.Seek(0, 0)
			_, fName := filepath.Split(newFile.Name())

			//set file meta
			fileMeta := meta.FileMeta{
				FileName: fName,
				Location: config.FileStoreDir + fName,
				FileSize: fileSize,
				FileSha1: fileHash,  //TODO file hash计算可能是个耗时的操作，后续考虑拆分
				UploadAt: time.Now(),
				UpdateAt: time.Now(),
				Status:   constant.FileStatusAvailable,
			}


			meta.UploadFileMetaDB(fileMeta)
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

func getFileMeta() gin.HandlerFunc {
	return func(context *gin.Context) {
		var req model.GetFileMetaReq

		if err := context.ShouldBind(&req); err != nil {
			log.Printf("bind request parameters error %v", err)
			context.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		fileMeta := meta.GetFileMetaDB(req.FileHash)
		if fileMeta.FileSha1 == "" {
			context.JSON(http.StatusOK, common.NewServiceResp(common.RespCodeSuccess, nil))
			return
		}
		context.JSON(http.StatusOK, common.NewServiceResp(common.RespCodeSuccess, fileMeta))
	}
}

func downloadHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		var req model.DownloadFileReq
		if err := context.ShouldBind(&req); err != nil {
			log.Printf("bind request parameters error %v", err)
			context.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		fm := meta.GetFileMetaDB(req.FileHash)
		if fm.FileSha1 == "" {
			context.JSON(http.StatusOK, common.NewServiceResp(common.RespCodeSuccess, nil))
			return
		}
		context.FileAttachment(fm.Location, fm.FileName)

	}
}

func updateFileMeta() gin.HandlerFunc {
	return func(context *gin.Context) {
		var req model.UpdateFileMetaReq
		if err := context.ShouldBind(&req); err != nil {
			log.Printf("bind request parameters error %v", err)
			context.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		fm := meta.GetFileMetaDB(req.FileHash)
		if fm.FileSha1 == "" {
			log.Printf("can't not found file meta %s", req.FileHash)
			context.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeNotFoundFileError, nil))
			return
		}

		if req.Filename != "" {
			fm.FileName = req.Filename
			meta.UpdateFileMetaDB(fm)

			//更新到user file关联表
			db.UpdateUserFilename(req.Username, req.FileHash, req.Filename)

			context.JSON(http.StatusOK,
				common.NewServiceResp(common.RespCodeSuccess, fm))
			return
		}

		log.Printf("filename cann't equals empty")
		context.JSON(http.StatusBadRequest,
			common.NewServiceResp(common.RespCodeFilenameError, nil))
	}
}

func deleteFile() gin.HandlerFunc {
	return func(context *gin.Context) {
		var req model.DeleteFileReq
		if err := context.ShouldBind(&req); err != nil {
			log.Printf("bind request parameters error %v", err)
			context.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		fm := meta.GetFileMetaDB(req.FileHash)
		if fm.FileSha1 == "" {
			log.Printf("can't not found file meta %s", req.FileHash)
			context.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeNotFoundFileError, nil))
			return
		}

		err := os.Remove(fm.Location)
		if err != nil {
			log.Printf("remove file error %v", err)
			context.JSON(http.StatusInternalServerError,
				common.NewServiceResp(common.RespCodeRemoveFileError, nil))
			return
		}

		meta.RemoveMetaDB(req.FileHash)
		context.JSON(http.StatusOK,
			common.NewServiceResp(common.RespCodeSuccess, nil))
	}
}

func queryFileList() gin.HandlerFunc {
	return func(context *gin.Context) {
		var req model.UserFileReq
		if err := context.ShouldBind(&req); err != nil {
			log.Printf("bind request parameters error %v", err)
			context.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		userFiles, ok := db.QueryUserFileMetas(req.Username, req.Limit)
		if !ok {
			context.JSON(http.StatusInternalServerError,
				common.NewServiceResp(common.RespCodeQueryFileError, nil))
			return
		}

		context.JSON(http.StatusOK,
			common.NewServiceResp(common.RespCodeSuccess, userFiles))
	}
}

func tryFastUpload() gin.HandlerFunc {
	return func(context *gin.Context) {
		var req model.FastUploadReq
		if err := context.ShouldBind(&req); err != nil {
			log.Printf("bind request parameters error %v", err)
			context.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		fm := meta.GetFileMetaDB(req.FileHash)
		if fm.FileSha1 == "" {
			log.Printf("not exist this file %s, failed to fast upload", req.FileHash)
			context.JSON(http.StatusOK,
				common.NewServiceResp(common.RespCodeFastUploadFailed, nil))
			return
		}

		ok := db.InsertUserFile(req.Username, req.FileHash, req.Filename, req.FileSize)
		if !ok {
			context.JSON(http.StatusInternalServerError,
				common.NewServiceResp(common.RespCodeUploadFileError, nil))
			return
		}

		context.JSON(http.StatusOK,
			common.NewServiceResp(common.RespCodeSuccess, nil))
	}
}