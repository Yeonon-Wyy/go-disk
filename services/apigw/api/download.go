package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	"go-disk/common/rpcinterface/downloadinterface"
	"go-disk/services/apigw/config"

	"log"
	"net/http"
)

var downloadCli downloadinterface.DownloadService

func init() {
	reg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			config.Conf.Micro.Registration.Consul.Addr,
		}
	})
	service := micro.NewService(
		micro.Registry(reg),
		micro.Name("go.micro.service.download"),
	)

	service.Init()

	downloadCli = downloadinterface.NewDownloadService("go.micro.service.download", service.Client())
}

func GetDownloadServiceEndpoint() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		resp, err :=downloadCli.DownloadEndPoint(context.TODO(), &downloadinterface.DownloadEndpointReq{})
		if err != nil {
			log.Printf("rpc call (get download service endpoint) error : %v", err)
			ctx.JSON(http.StatusBadRequest, *resp)
			return
		}

		ctx.JSON(http.StatusOK, *resp)
	}
}
