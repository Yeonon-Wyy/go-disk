package midware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	redisconn "go-disk/cache/redis"
	"go-disk/common"
	"go-disk/config"
	"go-disk/utils"
	"log"
	"net/http"
	"time"
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
	redisClient, err := redisconn.AuthConn()
	if err != nil {
		log.Printf("failed to connect redis server : %v", err)
		return ""
	}
	defer redisClient.Close()
	redisClient.Set(key, token, config.AuthRedisTokenExpireTime)
	return token
}

func ExistToken(key string) bool {
	redisClient, err := redisconn.AuthConn()
	if err != nil {
		log.Printf("failed to connect redis server : %v", err)
		return false
	}
	defer redisClient.Close()
	val := redisClient.Exists(key).Val()
	return val > 0
}

func ValidateToken(key, ReqToken string) bool {

	if !ExistToken(key) {
		log.Printf("the user not login!")
		return false
	}
	redisClient, err := redisconn.AuthConn()
	if err != nil {
		log.Printf("failed to connect redis server : %v", err)
		return false
	}
	defer redisClient.Close()

	internToken := redisClient.Get(key).Val()

	if internToken != ReqToken {
		log.Printf("token error!")
		return false
	}
	return true
}
