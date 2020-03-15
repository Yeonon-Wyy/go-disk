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
	authCli authinterface.AuthService
	downloadCli downloadinterface.DownloadService
	fileCli fileinterface.FileService
	uploadCli uploadinterface.UploadService
	userCli userinterface.UserService
)

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

func GetDownloadCli() downloadinterface.DownloadService {
	if downloadCli != nil {
		return downloadCli
	}
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

	downloadCli = downloadinterface.NewDownloadService(config.Conf.Micro.Client.Download.ServiceName, serv.Client())
	return downloadCli
}

func GetFileCli() fileinterface.FileService {
	if fileCli != nil {
		return fileCli
	}
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

	fileCli = fileinterface.NewFileService("go.micro.service.file", serv.Client())
	return fileCli
}

func GetUploadCli() uploadinterface.UploadService {
	if uploadCli != nil {
		return uploadCli
	}

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

	uploadCli = uploadinterface.NewUploadService(config.Conf.Micro.Client.Upload.ServiceName, serv.Client())
	return uploadCli
}

func GetUserCli() userinterface.UserService {
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

	userCli = userinterface.NewUserService(config.Conf.Micro.Client.User.ServiceName, service.Client())
	return userCli
}
