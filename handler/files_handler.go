package handler

import (
	"github.com/gin-gonic/gin"
	"go-disk/common"
	"go-disk/config"
	"go-disk/meta"
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
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		meta.UpdateFileMeta(fileMeta.FileSha1, fileMeta)

		log.Printf("upload file success, file hash is : %s", fileMeta.FileSha1)

		context.JSON(http.StatusOK, common.NewServiceResp(common.RespCodeSuccess, nil))
	}
}

func getFileMeta() gin.HandlerFunc {
	return func(context *gin.Context) {
		var req GetFileMetaReq

		if err := context.ShouldBind(&req); err != nil {
			log.Printf("bind request parameters error %v", err)
			context.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		fileMeta := meta.GetFileMeta(req.FileHash)
		context.JSON(http.StatusOK, common.NewServiceResp(common.RespCodeSuccess, fileMeta))
	}
}

func downloadHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		var req DownloadFileReq
		if err := context.ShouldBind(&req); err != nil {
			log.Printf("bind request parameters error %v", err)
			context.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		fm := meta.GetFileMeta(req.FileHash)

		context.FileAttachment(fm.Location, fm.FileName)

	}
}

func updateFileMeta() gin.HandlerFunc {
	return func(context *gin.Context) {
		var req UpdateFileMetaReq
		if err := context.ShouldBind(&req); err != nil {
			log.Printf("bind request parameters error %v", err)
			context.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		if meta.GetFileMeta(req.FileHash) == nil {
			log.Printf("can't not found file meta %s", req.FileHash)
			context.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		if req.Filename != "" {
			fm := meta.GetFileMeta(req.FileHash)
			fm.FileName = req.Filename
			meta.UpdateFileMeta(req.FileHash, *fm)
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
		var req DeleteFileReq
		if err := context.ShouldBind(&req); err != nil {
			log.Printf("bind request parameters error %v", err)
			context.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		fm := meta.GetFileMeta(req.FileHash)
		if fm == nil {
			log.Printf("can't not found file meta %s", req.FileHash)
			context.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		err := os.Remove(fm.Location)
		if err != nil {
			log.Printf("remove file error %v", err)
			context.JSON(http.StatusInternalServerError,
				common.NewServiceResp(common.RespCodeRemoveFileError, nil))
			return
		}

		meta.RemoveMeta(req.FileHash)
		context.JSON(http.StatusOK,
			common.NewServiceResp(common.RespCodeSuccess, nil))


	}
}