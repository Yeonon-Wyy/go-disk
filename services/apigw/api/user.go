package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-disk/common"
	"go-disk/common/log4disk"
	"go-disk/common/rpcinterface/userinterface"
	"go-disk/services/apigw/rpc"
	"go-disk/services/apigw/vo"
	"net/http"
)

func RegisterUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req vo.UserRegisterReq
		if err := ctx.ShouldBind(&req); err != nil {
			log4disk.E("request parameters error : %v", err)
			ctx.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		resp, err := rpc.UserCli.UserRegister(context.TODO(), &userinterface.RegisterReq{
			Username: req.Username,
			Password: req.Password,
		})

		if err != nil || resp.Code != int64(common.RespCodeSuccess.Code) {
			log4disk.E("rpc call (user register) error : %v", err)
			ctx.JSON(http.StatusInternalServerError, *resp)
			return
		}

		ctx.JSON(http.StatusOK, *resp)
	}
}

func QueryUserInfo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req vo.UserQueryReq
		if err := ctx.ShouldBindUri(&req); err != nil {
			log4disk.E("request parameters error : %v", err)
			ctx.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		resp, err := rpc.UserCli.QueryUserInfo(context.TODO(), &userinterface.QueryUserInfoReq{
			Username: req.Username,
		})

		if err != nil || resp.Code != int64(common.RespCodeSuccess.Code) {
			log4disk.E("rpc call (query user info) error : %v", err)
			ctx.JSON(http.StatusInternalServerError, *resp)
			return
		}

		ctx.JSON(http.StatusOK, *resp)
	}
}
