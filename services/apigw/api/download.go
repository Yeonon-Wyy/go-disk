package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-disk/common/log4disk"
	"go-disk/common/rpcinterface/downloadinterface"
	"go-disk/services/apigw/rpc"

	"net/http"
)

func GetDownloadServiceEndpoint() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		resp, err := rpc.DownloadCli.DownloadEndPoint(context.TODO(), &downloadinterface.DownloadEndpointReq{})
		if err != nil {
			log4disk.E("rpc call (get download service endpoint) error : %v", err)
			ctx.JSON(http.StatusInternalServerError, *resp)
			return
		}

		ctx.JSON(http.StatusOK, *resp)
	}
}
