package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-disk/common"
	"go-disk/utils"
	"log"
	"net/http"
	"time"
)



var (
	//TODO: just temp
	tokenCache = make(map[string]string)
)

func AuthorizeInterceptor() gin.HandlerFunc {
	return func(context *gin.Context) {

		username, token := context.Query("username"), context.Query("token")

		if username == "" || token == "" {
			log.Printf("request param error")
			context.Abort()
			context.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		if !ValidateToken(username, token) {
			log.Printf("token validate error")
			context.Abort()
			context.JSON(http.StatusUnauthorized,
				common.NewServiceResp(common.RespCodeUnauthorizedError, nil))
			return
		}

		context.Next()
	}
}

func GenToken(key string) string {
	ts := fmt.Sprintf("%x", time.Now().Unix())
	token := utils.MD5([]byte(key + ts + "_tokensalt")) + ts[:8]
	tokenCache[key] = token
	return token
}

func ExistToken(key string) bool {
	 _, ok := tokenCache[key]
	 return ok
}

//TODO: Token应该是有时效性的，这里暂时不考虑，后续用redis实现
func ValidateToken(key, ReqToken string) bool {
	internToken, ok := tokenCache[key]
	if !ok {
		log.Printf("the user not login!")
		return false
	}

	if internToken != ReqToken {
		log.Printf("token error!")
		return false
	}
	return true
}
