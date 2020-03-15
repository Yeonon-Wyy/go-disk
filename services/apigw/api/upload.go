package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-disk/common/rpcinterface/uploadinterface"
	"go-disk/services/apigw/rpc"
	"log"
	"net/http"
)


func GetUploadServiceEndpoint() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		resp, err := rpc.GetUploadCli().UploadEndPoint(context.TODO(), &uploadinterface.UploadEndPointReq{})
		if err != nil {
			log.Printf("rpc call (get upload service endpoint) error : %v", err)
			ctx.JSON(http.StatusInternalServerError, *resp)
			return
		}

		ctx.JSON(http.StatusOK, *resp)
	}
}