package rpc

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	"go-disk/common/rpcinterface/authinterface"
	"go-disk/services/upload/config"
)

var authCli authinterface.AuthService

func GetAuthCli() authinterface.AuthService {
	if authCli != nil {
		return authCli
	}
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

	authCli = authinterface.NewAuthService(config.Conf.Micro.Client.Auth.ServiceName, serv.Client())
	return authCli
}
