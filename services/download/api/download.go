package api

import (
	"github.com/gin-gonic/gin"
	"go-disk/common"
	"go-disk/config"
	"go-disk/meta"
	"go-disk/store/ceph"
	"net/http"
	"strconv"
)

func DownloadHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		//var req model.DownloadFileReq
		//if err := context.ShouldBind(&req); err != nil {
		//	log.Printf("bind request parameters error %v", err)
		//	context.JSON(http.StatusBadRequest,
		//		common.NewServiceResp(common.RespCodeBindReParamError, nil))
		//	return
		//}
		fileHash := context.Query("file_hash")

		fm := meta.GetFileMetaDB(fileHash)
		if fm.FileSha1 == "" {
			context.JSON(http.StatusOK, common.NewServiceResp(common.RespCodeSuccess, nil))
			return
		}

		bucket := ceph.GetCephBucket(config.CephFileStoreBucketName)
		fileData, _ := bucket.Get(fm.Location)

		context.Writer.Header().Set("Content-Type", "application/octet-stream")
		context.Writer.Header().Set("Content-Disposition", "attachment; filename=\""+fm.FileName+"\"")
		context.Writer.Header().Set("Content-Length", strconv.Itoa(len(fileData)))
		context.Writer.Write(fileData)

	}
}