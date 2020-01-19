package handler

import (
	"github.com/gin-gonic/gin"
	"go-disk/common"
	"go-disk/config"
	"go-disk/meta"
	"go-disk/model"
	"go-disk/utils"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type FilesServiceHandler struct {
	BashPath string
}


func (f FilesServiceHandler) Init(group *gin.RouterGroup) {
	group.POST("/upload", uploadFile())

	group.GET("/meta", getFileMeta())
	group.GET("/download", downloadHandler())
	group.PUT("/meta", updateFileMeta())
	group.DELETE("/delete", deleteFile())
}

func uploadFile() gin.HandlerFunc {
	return func(context *gin.Context) {
		fh, err := context.FormFile("file")
		if err != nil {
			log.Printf("read from file error : %v", err)
			context.JSON(http.StatusInternalServerError, common.NewServiceResp(common.RespCodeReadFileError, nil))
			return
		}

		newFile, err := os.Create(config.FileStoreDir + fh.Filename)
		if err != nil {
			log.Printf("create file error : %v", err)
			context.JSON(http.StatusInternalServerError, common.NewServiceResp(common.RespCodeCreateFileError, nil))
			return
		}

		file, err := fh.Open()
		if err != nil {
			log.Printf("open file error : %v", err)
			context.JSON(http.StatusInternalServerError, common.NewServiceResp(common.RespCodeOpenFileError, nil))
			return
		}
		defer file.Close()

		fileSize, err := io.Copy(newFile, file)

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
			FileSha1: utils.FileSha1(newFile),
			UploadAt: time.Now(),
			UpdateAt: time.Now(),
			Status: common.FileStatusAvailable,
		}


		meta.UploadFileMetaDB(fileMeta)

		log.Printf("upload file success, file hash is : %s", fileMeta.FileSha1)

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