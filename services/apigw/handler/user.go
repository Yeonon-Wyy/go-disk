package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	"go-disk/common"
	"go-disk/config"
	"go-disk/model"
	userproto "go-disk/services/user/proto"
	"log"
	"net/http"
)

var userCli userproto.UserService

func init() {

	reg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			config.ConsulAddress,
		}
	})
	service := micro.NewService(
		micro.Registry(reg),
		micro.Name("go.micro.service.user"),
		)

	service.Init()

	userCli = userproto.NewUserService("go.micro.service.user", service.Client())
}

func RegisterUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req model.UserRegisterReq
		if err := ctx.ShouldBind(&req); err != nil {
			log.Printf("request parameters error : %v", err)
			ctx.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		resp, err := userCli.UserRegister(context.TODO(), &userproto.RegisterReq{
			Username:             req.Username,
			Password:             req.Password,
		})

		if err != nil {
			log.Printf("rpc call (user register) error : %v", err)
			ctx.JSON(http.StatusBadRequest, *resp)
			return
		}

		ctx.JSON(http.StatusOK, *resp)
	}
}

func UserLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req model.UserLoginReq
		if err := ctx.ShouldBind(&req); err != nil {
			log.Printf("request parameters error : %v", err)
			ctx.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		resp, err := userCli.UserLogin(context.TODO(), &userproto.LoginReq{
			Username: req.Username,
			Password: req.Password,
		})

		if err != nil {
			log.Printf("rpc call (user login) error : %v", err)
			ctx.JSON(http.StatusInternalServerError, *resp)
			return
		}

		ctx.JSON(http.StatusOK,*resp)
	}
}

func QueryUserInfo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req model.UserQueryReq
		if err := ctx.ShouldBind(&req); err != nil {
			log.Printf("request parameters error : %v", err)
			ctx.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		resp, err := userCli.QueryUserInfo(context.TODO(), &userproto.QueryUserInfoReq{
			Username: req.Username,
			AccessToken: req.Token,
		})

		if err != nil {
			log.Printf("rpc call (query user info) error : %v", err)
			ctx.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		ctx.JSON(http.StatusOK,*resp)
	}
}
