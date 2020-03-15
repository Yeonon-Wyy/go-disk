package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-disk/common"
	"go-disk/common/rpcinterface/authinterface"
	"go-disk/common/rpcinterface/userinterface"
	"go-disk/services/apigw/rpc"
	"go-disk/services/apigw/vo"
	"log"
	"net/http"
)



func RegisterUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req vo.UserRegisterReq
		if err := ctx.ShouldBind(&req); err != nil {
			log.Printf("request parameters error : %v", err)
			ctx.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		resp, err := rpc.GetUserCli().UserRegister(context.TODO(), &userinterface.RegisterReq{
			Username:             req.Username,
			Password:             req.Password,
		})

		if err != nil || resp.Code != int64(common.RespCodeSuccess.Code) {
			log.Printf("rpc call (user register) error : %v", err)
			ctx.JSON(http.StatusBadRequest, *resp)
			return
		}

		ctx.JSON(http.StatusOK, *resp)
	}
}

func UserLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req vo.UserLoginReq
		if err := ctx.ShouldBind(&req); err != nil {
			log.Printf("request parameters error : %v", err)
			ctx.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		resp, err := rpc.GetAuthCli().Authorize(context.TODO(), &authinterface.AuthorizeReq{
			Username: req.Username,
			Password: req.Password,
		})


		if err != nil || resp.Code != int64(common.RespCodeSuccess.Code) {
			log.Printf("rpc call (user login) error : %v", err)
			ctx.JSON(http.StatusInternalServerError, *resp)
			return
		}

		ctx.JSON(http.StatusOK, *resp)
	}
}

func QueryUserInfo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req vo.UserQueryReq
		if err := ctx.ShouldBind(&req); err != nil {
			log.Printf("request parameters error : %v", err)
			ctx.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		resp, err := rpc.GetUserCli().QueryUserInfo(context.TODO(), &userinterface.QueryUserInfoReq{
			Username: req.Username,
			AccessToken: req.Token,
		})

		if err != nil || resp.Code != int64(common.RespCodeSuccess.Code) {
			log.Printf("rpc call (query user info) error : %v", err)
			ctx.JSON(http.StatusBadRequest, *resp)
			return
		}

		ctx.JSON(http.StatusOK,*resp)
	}
}
