package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-disk/common"
	"go-disk/common/log4disk"
	"go-disk/common/rpcinterface/authinterface"
	"go-disk/services/apigw/rpc"
	"go-disk/services/apigw/vo"
	"net/http"
)

func Authorize() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req vo.AuthorizeReq
		if err := ctx.ShouldBind(&req); err != nil {
			log4disk.E("request parameters error : %v", err)
			ctx.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		resp, err := rpc.AuthCli.Authorize(context.TODO(), &authinterface.AuthorizeReq{
			Username: req.Username,
			Password: req.Password,
		})

		if err != nil || resp.Code != int64(common.RespCodeSuccess.Code) {
			log4disk.E("rpc call (user login) error : %v", err)
			ctx.JSON(http.StatusInternalServerError, *resp)
			return
		}

		ctx.JSON(http.StatusOK, *resp)
	}
}

func UnAuthorize() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req vo.UnAuthorizeReq
		if err := ctx.ShouldBindUri(&req); err != nil {
			log4disk.E("request parameters error : %v", err)
			ctx.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		resp, err := rpc.AuthCli.UnAuthorize(context.TODO(), &authinterface.UnAuthorizeReq{
			Username: req.Username,
		})

		if err != nil || resp.Code != int64(common.RespCodeSuccess.Code) {
			log4disk.E("rpc call (user login) error : %v", err)
			ctx.JSON(http.StatusInternalServerError, *resp)
			return
		}

		ctx.JSON(http.StatusOK, *resp)
	}
}
