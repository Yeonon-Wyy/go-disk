package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	"go-disk/config"
	uploadproto "go-disk/services/upload/proto"
	"log"
	"net/http"
)

var uploadCli uploadproto.UploadService

func init() {

	reg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			config.ConsulAddress,
		}
	})
	service := micro.NewService(
		micro.Registry(reg),
		micro.Name("go.micro.service.upload"),
	)

	service.Init()

	uploadCli = uploadproto.NewUploadService("go.micro.service.upload", service.Client())
}

func GetUploadServiceEndpoint() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		resp, err := uploadCli.UploadEndPoint(context.TODO(), &uploadproto.UploadEndPointReq{})
		if err != nil {
			log.Printf("rpc call (get upload service endpoint) error : %v", err)
			ctx.JSON(http.StatusBadRequest, *resp)
			return
		}

		ctx.JSON(http.StatusOK, *resp)
	}
}