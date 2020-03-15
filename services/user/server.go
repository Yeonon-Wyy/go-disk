package main

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	"go-disk/common/rpcinterface/userinterface"
	"go-disk/services/user/config"
	"go-disk/services/user/rpc"
)

func main() {
	reg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			config.Conf.Micro.Registration.Consul.Addr,
		}
	})
	service := micro.NewService(
		micro.Registry(reg),
		micro.Name("go.micro.service.user"))

	service.Init()

	err := userinterface.RegisterUserServiceHandler(service.Server(), new(rpc.UserHandler))
	if err != nil {
		panic(err)
	}

	err = service.Run()
	if err != nil {
		panic(err)
	}
}
