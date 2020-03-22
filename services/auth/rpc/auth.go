package rpc

import (
	"context"
	"go-disk/common"
	"go-disk/common/jwt"
	"go-disk/common/log4disk"
	"go-disk/common/rpcinterface/authinterface"
	"go-disk/common/utils"
	"go-disk/services/auth/config"
	"go-disk/services/auth/db"
	redisconn "go-disk/services/auth/redis"
	"time"
)

type AuthServiceHandler struct {
}

func (a *AuthServiceHandler) Authentication(ctx context.Context, req *authinterface.AuthenticationReq, resp *authinterface.AuthenticationResp) error {
	success := validateToken(req.AccessToken)
	if !success {
		resp.Code = int64(common.RespCodeValidateTokenError.Code)
		resp.Message = common.RespCodeValidateTokenError.Message
		return nil
	}
	resp.Code = int64(common.RespCodeSuccess.Code)
	resp.Message = common.RespCodeSuccess.Message

	return nil
}

func (a *AuthServiceHandler) Authorize(ctx context.Context, req *authinterface.AuthorizeReq, resp *authinterface.AuthorizeResp) error {
	sha1Pwd := utils.Sha1([]byte(req.Password + config.Conf.Business.UserPasswordSalt))
	userExist := db.ExistUserByUsernameAndPassword(req.Username, sha1Pwd)
	if !userExist {
		resp.Code = int64(common.RespCodeUserNotFound.Code)
		resp.Message = common.RespCodeUserNotFound.Message
		return nil
	}

	tokenStr, err := genToken(req.Username)
	if err != nil {
		resp.Code = int64(common.RespCodeGenTokenError.Code)
		resp.Message = common.RespCodeGenTokenError.Message
		return nil
	}

	resp.Code = int64(common.RespCodeSuccess.Code)
	resp.Message = common.RespCodeSuccess.Message
	resp.Data = &authinterface.AuthorizeResp_Data{
		AccessToken: tokenStr,
	}

	return nil
}

func (a *AuthServiceHandler) UnAuthorize(ctx context.Context, req *authinterface.UnAuthorizeReq, resp *authinterface.UnAuthorizeResp) error {
	if ok := deleteToken(req.Username); ok {
		resp.Code = int64(common.RespCodeSuccess.Code)
		resp.Message = common.RespCodeSuccess.Message
	} else {
		resp.Code = int64(common.RespCodeDeleteTokenError.Code)
		resp.Message = common.RespCodeDeleteTokenError.Message
	}
	return nil
}

func genToken(username string) (string, error) {

	tokenString, err := jwt.GenToken(map[string]interface{}{
		"username": username,
	})

	if err != nil {
		return "", err
	}

	redisClient, err := redisconn.AuthConn()
	defer redisClient.Close()
	if err != nil {
		log4disk.E("failed to connect redis server : %v", err)
		return "", err
	}

	expireTime := time.Duration(config.Conf.DataSource.Redis.TokenExpireTime) * time.Hour
	redisClient.Set(username, tokenString, expireTime)
	return tokenString, nil
}

func existToken(username string) bool {
	redisClient, err := redisconn.AuthConn()
	if err != nil {
		log4disk.E("failed to connect redis server : %v", err)
		return false
	}
	defer redisClient.Close()
	val := redisClient.Exists(username).Val()
	return val > 0
}

func validateToken(ReqTokenStr string) bool {

	payload, suc := jwt.GetPayload(ReqTokenStr)
	if !suc {
		return false
	}

	username := payload["username"].(string)

	if !existToken(username) {
		log4disk.E("the user not login!")
		return false
	}
	return true

}

func deleteToken(username string) bool {
	redisClient, err := redisconn.AuthConn()
	if err != nil {
		log4disk.E("failed to connect redis server : %v", err)
		return false
	}
	defer redisClient.Close()

	res := redisClient.Del(username).Val()
	return res > 0
}
