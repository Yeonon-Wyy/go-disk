package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	"go-disk/common/rpcinterface/uploadinterface"
	"go-disk/services/apigw/config"
	"log"
	"net/http"
)

var uploadCli uploadinterface.UploadService

func init() {

	reg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			config.Conf.Micro.Registration.Consul.Addr,
		}
	})
	service := micro.NewService(
		micro.Registry(reg),
		micro.Name("go.micro.service.upload"),
	)

	service.Init()

	uploadCli = uploadinterface.NewUploadService("go.micro.service.upload", service.Client())
}

func GetUploadServiceEndpoint() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		resp, err := uploadCli.UploadEndPoint(context.TODO(), &uploadinterface.UploadEndPointReq{})
		if err != nil {
			log.Printf("rpc call (get upload service endpoint) error : %v", err)
			ctx.JSON(http.StatusBadRequest, *resp)
			return
		}

		ctx.JSON(http.StatusOK, *resp)
	}
}