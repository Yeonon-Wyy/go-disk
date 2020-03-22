package rpc

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	"go-disk/common/rpcinterface/authinterface"
	"go-disk/services/download/config"
)

var AuthCli authinterface.AuthService

func init() {
	initAuthCli()
}

func initAuthCli() {
	reg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			config.Conf.Micro.Registration.Consul.Addr,
		}
	})
	serv := micro.NewService(
		micro.Registry(reg),
		micro.Name(config.Conf.Micro.Client.Auth.ServiceName),
	)

	serv.Init()

	AuthCli = authinterface.NewAuthService(config.Conf.Micro.Client.Auth.ServiceName, serv.Client())
}
