package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-disk/common/rpcinterface/downloadinterface"
	"go-disk/services/apigw/rpc"

	"log"
	"net/http"
)


func GetDownloadServiceEndpoint() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		resp, err := rpc.GetDownloadCli().DownloadEndPoint(context.TODO(), &downloadinterface.DownloadEndpointReq{})
		if err != nil {
			log.Printf("rpc call (get download service endpoint) error : %v", err)
			ctx.JSON(http.StatusBadRequest, *resp)
			return
		}

		ctx.JSON(http.StatusOK, *resp)
	}
}
