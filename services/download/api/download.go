package api

import (
	"github.com/gin-gonic/gin"
	"go-disk/common"
	"go-disk/config"
	"go-disk/db"
	"go-disk/store/ceph"
	"net/http"
	"strconv"
)

func DownloadHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		fileHash := context.Query("file_hash")

		tblFile, err := db.GetFileMeta(fileHash)
		if err != nil {
			context.JSON(http.StatusOK, common.NewServiceResp(common.RespCodeSuccess, nil))
			return
		}

		bucket := ceph.GetCephBucket(config.CephFileStoreBucketName)
		fileData, _ := bucket.Get(tblFile.FileLocation)

		context.Writer.Header().Set("Content-Type", "application/octet-stream")
		context.Writer.Header().Set("Content-Disposition", "attachment; filename=\""+tblFile.Filename+"\"")
		context.Writer.Header().Set("Content-Length", strconv.Itoa(len(fileData)))
		context.Writer.Write(fileData)

	}
}