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
	"time"
)

type FilesServiceHandler struct {
	BashPath string
}

func (f FilesServiceHandler) Init(group *gin.RouterGroup) {
	group.POST("/upload", uploadFile())

	group.GET("/meta", getFileMeta())
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

		//set file meta
		fileMeta := meta.FileMeta{
			FileName: newFile.Name(),
			Location: "/Volumes/computer/go/go-disk/filestore/" + newFile.Name(),
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
		fileHash := context.Query("file_hash")
		fileMeta := meta.GetFileMeta(fileHash)

		context.JSON(http.StatusOK, common.NewServiceResp(common.RespCodeSuccess, fileMeta))
	}
}