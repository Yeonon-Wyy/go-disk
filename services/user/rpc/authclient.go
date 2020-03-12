package rpc

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	"go-disk/common/rpcinterface/authinterface"
	"go-disk/services/user/config"
)

var AuthCli authinterface.AuthService

func init() {
	reg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			config.Conf.Micro.Registration.Consul.Addr,
		}
	})
	service := micro.NewService(
		micro.Registry(reg),
		micro.Name("go.micro.service.auth"),
	)

	service.Init()

	AuthCli = authinterface.NewAuthService("go.micro.service.auth", service.Client())
}
