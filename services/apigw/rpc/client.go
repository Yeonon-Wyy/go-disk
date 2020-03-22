package rpc

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	"go-disk/common/rpcinterface/authinterface"
	"go-disk/common/rpcinterface/downloadinterface"
	"go-disk/common/rpcinterface/fileinterface"
	"go-disk/common/rpcinterface/uploadinterface"
	"go-disk/common/rpcinterface/userinterface"
	"go-disk/services/apigw/config"
)

var (
	AuthCli     authinterface.AuthService
	DownloadCli downloadinterface.DownloadService
	FileCli     fileinterface.FileService
	UploadCli   uploadinterface.UploadService
	UserCli     userinterface.UserService
)

func init() {
	initAuthCli()
	initDownloadCli()
	initFileCli()
	initUploadCli()
	initUserCli()
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

func initDownloadCli() {
	reg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			config.Conf.Micro.Registration.Consul.Addr,
		}
	})
	serv := micro.NewService(
		micro.Registry(reg),
		micro.Name(config.Conf.Micro.Client.Download.ServiceName),
	)

	serv.Init()

	DownloadCli = downloadinterface.NewDownloadService(config.Conf.Micro.Client.Download.ServiceName, serv.Client())
}

func initFileCli() {
	reg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			config.Conf.Micro.Registration.Consul.Addr,
		}
	})
	serv := micro.NewService(
		micro.Registry(reg),
		micro.Name("go.micro.service.file"),
	)

	serv.Init()

	FileCli = fileinterface.NewFileService("go.micro.service.file", serv.Client())
}

func initUploadCli() {
	reg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			config.Conf.Micro.Registration.Consul.Addr,
		}
	})
	serv := micro.NewService(
		micro.Registry(reg),
		micro.Name(config.Conf.Micro.Client.Upload.ServiceName),
	)

	serv.Init()

	UploadCli = uploadinterface.NewUploadService(config.Conf.Micro.Client.Upload.ServiceName, serv.Client())
}

func initUserCli() {
	reg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			config.Conf.Micro.Registration.Consul.Addr,
		}
	})
	service := micro.NewService(
		micro.Registry(reg),
		micro.Name(config.Conf.Micro.Client.User.ServiceName),
	)

	service.Init()

	UserCli = userinterface.NewUserService(config.Conf.Micro.Client.User.ServiceName, service.Client())
}
