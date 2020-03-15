package main

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	"go-disk/common/rpcinterface/fileinterface"
	"go-disk/services/file/config"
	"go-disk/services/file/rpc"
)

func main() {
	reg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			config.Conf.Micro.Registration.Consul.Addr,
		}
	})
	service := micro.NewService(
		micro.Registry(reg),
		micro.Name("go.micro.service.file"))

	service.Init()

	err := fileinterface.RegisterFileServiceHandler(service.Server(), new(rpc.FileService))
	if err != nil {
		panic(err)
	}

	err = service.Run()
	if err != nil {
		panic(err)
	}
}
