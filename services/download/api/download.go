package api

import (
	"github.com/gin-gonic/gin"
	"go-disk/common"
	"go-disk/services/download/config"
	"go-disk/services/download/db"
	"go-disk/services/download/store/ceph"
	"net/http"
	"strconv"
)

func DownloadHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		fileHash := context.Query("file_hash")

		tblFile, err := db.GetFileMeta(fileHash)
		if err != nil {
			context.JSON(http.StatusInternalServerError, common.NewServiceResp(common.RespCodeNotFoundFileError, nil))
			return
		}

		bucket := ceph.GetCephBucket(config.Conf.Store.Ceph.FileStoreBucketName)
		fileData, _ := bucket.Get(tblFile.FileAddr)

		context.Writer.Header().Set("Content-Type", "application/octet-stream")
		context.Writer.Header().Set("Content-Disposition", "attachment; filename=\""+tblFile.FileName+"\"")
		context.Writer.Header().Set("Content-Length", strconv.Itoa(len(fileData)))
		context.Writer.Write(fileData)

	}
}