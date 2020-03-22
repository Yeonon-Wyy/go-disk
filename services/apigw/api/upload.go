package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-disk/common/log4disk"
	"go-disk/common/rpcinterface/uploadinterface"
	"go-disk/services/apigw/rpc"
	"net/http"
)

func GetUploadServiceEndpoint() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		resp, err := rpc.UploadCli.UploadEndPoint(context.TODO(), &uploadinterface.UploadEndPointReq{})
		if err != nil {
			log4disk.E("rpc call (get upload service endpoint) error : %v", err)
			ctx.JSON(http.StatusInternalServerError, *resp)
			return
		}

		ctx.JSON(http.StatusOK, *resp)
	}
}
