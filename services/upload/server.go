package main

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	"go-disk/common/log4disk"
	"go-disk/common/rpcinterface/uploadinterface"
	"go-disk/services/upload/config"
	"go-disk/services/upload/router"
	"go-disk/services/upload/rpc"
)

func main() {
	uploadRouter := router.Router()

	go startRpcService()

	err := uploadRouter.Run(":9000")
	if err != nil {
		panic(err)
	}
}

func startRpcService() {
	reg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			config.Conf.Micro.Registration.Consul.Addr,
		}
	})
	service := micro.NewService(
		micro.Registry(reg),
		micro.Name("go.micro.service.upload"))

	service.Init()

	err := uploadinterface.RegisterUploadServiceHandler(service.Server(), new(rpc.EndPoint))
	if err != nil {
		panic(err)
	}

	err = service.Run()
	if err != nil {
		panic(err)
	}

	log4disk.I("start upload service success")
}
