package interceptor

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-disk/common"
	"go-disk/common/rpcinterface/authinterface"
	"go-disk/services/apigw/rpc"
	"log"
	"net/http"
)

func AuthorizeInterceptor() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")

		if token == "" {
			log.Printf("request param error")
			ctx.Abort()
			ctx.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}
		resp, err := rpc.AuthCli.Authentication(context.TODO(), &authinterface.AuthenticationReq{
			AccessToken: token,
		})
		if err != nil || resp.Code != int64(common.RespCodeSuccess.Code) {
			log.Printf("token validate error : %v", err)
			ctx.Abort()
			ctx.JSON(http.StatusUnauthorized,
				common.NewServiceResp(common.RespCodeUnauthorizedError, nil))
			return
		}

		ctx.Next()
	}
}


