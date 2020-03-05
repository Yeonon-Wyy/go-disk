package main

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	"go-disk/common/rpcinterface/fileinterface"
	"go-disk/config"
	"go-disk/services/file/handler"
)

func main() {
	reg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			config.ConsulAddress,
		}
	})
	service := micro.NewService(
		micro.Registry(reg),
		micro.Name("go.micro.service.file"))

	service.Init()

	err := fileinterface.RegisterFileServiceHandler(service.Server(), new(handler.FileService))
	if err != nil {
		panic(err)
	}

	err = service.Run()
	if err != nil {
		panic(err)
	}
}
