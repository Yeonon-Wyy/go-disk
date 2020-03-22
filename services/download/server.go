package main

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	"go-disk/common/log4disk"
	"go-disk/common/rpcinterface/downloadinterface"
	"go-disk/services/download/config"
	"go-disk/services/download/router"
	"go-disk/services/download/rpc"
)

func main() {
	route := router.Router()

	go startRpcService()

	err := route.Run(":8000")
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
		micro.Name("go.micro.service.download"))

	service.Init()

	err := downloadinterface.RegisterDownloadServiceHandler(service.Server(), new(rpc.Download))
	if err != nil {
		panic(err)
	}

	err = service.Run()
	if err != nil {
		panic(err)
	}

	log4disk.I("start upload service success")
}