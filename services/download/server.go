package main

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	"go-disk/config"
	"go-disk/services/download/proto"
	"go-disk/services/download/router"
	"go-disk/services/download/rpc"
	"log"
)

func main() {
	route := router.Router()

	go startRpcService()

	err := route.Run()
	if err != nil {
		panic(err)
	}
}

func startRpcService() {
	reg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			config.ConsulAddress,
		}
	})
	service := micro.NewService(
		micro.Registry(reg),
		micro.Name("go.micro.service.download"))

	service.Init()

	err := proto.RegisterDownloadServiceHandler(service.Server(), new(rpc.Download))
	if err != nil {
		panic(err)
	}

	err = service.Run()
	if err != nil {
		panic(err)
	}

	log.Printf("start upload service success")
}