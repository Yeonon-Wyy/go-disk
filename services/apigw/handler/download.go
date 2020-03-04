package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	"go-disk/config"
	downloadProto "go-disk/services/download/proto"
	"log"
	"net/http"
)

var downloadCli downloadProto.DownloadService

func init() {
	reg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			config.ConsulAddress,
		}
	})
	service := micro.NewService(
		micro.Registry(reg),
		micro.Name("go.micro.service.download"),
	)

	service.Init()

	downloadCli = downloadProto.NewDownloadService("go.micro.service.download", service.Client())
}

func GetDownloadServiceEndpoint() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		resp, err :=downloadCli.DownloadEndPoint(context.TODO(), &downloadProto.DownloadEndpointReq{})
		if err != nil {
			log.Printf("rpc call (get download service endpoint) error : %v", err)
			ctx.JSON(http.StatusBadRequest, *resp)
			return
		}

		ctx.JSON(http.StatusOK, *resp)
	}
}
