package handler

import (
	"context"
	"errors"
	"go-disk/common"
	userrpcinterface "go-disk/common/rpcinterface/userinterface"
	"go-disk/services/apigw/auth"
	"go-disk/services/user/config"
	"go-disk/services/user/db"

	"go-disk/utils"
)

type UserHandler struct {

}


func (u UserHandler) UserRegister(ctx context.Context, req *userrpcinterface.RegisterReq, resp *userrpcinterface.RegisterResp) error {

	if db.ExistUserByUsername(req.Username) {
		resp.Code = int64(common.RespCodeUserAlreadyRegistered.Code)
		resp.Message = common.RespCodeUserAlreadyRegistered.Message
		return errors.New("user already register")
	}

	pwd := utils.Sha1([]byte(req.Password + config.PwdSalt))

	if !db.InsertUser(req.Username, pwd) {
		resp.Code = int64(common.RespCodeUserRegisterError.Code)
		resp.Message = common.RespCodeUserRegisterError.Message
		return errors.New("user register error")
	}

	resp.Code = int64(common.RespCodeSuccess.Code)
	resp.Message = common.RespCodeSuccess.Message
	return nil
}

func (u UserHandler) UserLogin(ctx context.Context, req *userrpcinterface.LoginReq, resp *userrpcinterface.LoginResp) error {
	if auth.ExistToken(req.Username) {
		resp.Code = int64(common.RespCodeUserAlreadyLogin.Code)
		resp.Message = common.RespCodeUserAlreadyLogin.Message
		return errors.New("user already login")
	}

	exist := db.ExistUserByUsernameAndPassword(req.Username, utils.Sha1([]byte(req.Password + config.PwdSalt)))
	if !exist {
		resp.Code = int64(common.RespCodeUserNotFound.Code)
		resp.Message = common.RespCodeUserNotFound.Message
		return errors.New("not this user")
	}

	resp.Code = int64(common.RespCodeSuccess.Code)
	resp.Message = common.RespCodeSuccess.Message
	resp.Data = &userrpcinterface.LoginResp_Data{AccessToken: auth.GenToken(req.Username)}

	return nil
}

func (u UserHandler) QueryUserInfo(ctx context.Context, req *userrpcinterface.QueryUserInfoReq, resp *userrpcinterface.QueryUserInfoResp) error {

	result, err := db.QueryUser(req.Username)
	if err != nil {
		resp.Code = int64(common.RespCodeUserNotFound.Code)
		resp.Message = common.RespCodeUserNotFound.Message
		return errors.New("user not found")
	}

	resp.Code = int64(common.RespCodeSuccess.Code)
	resp.Message = common.RespCodeSuccess.Message
	resp.Data = &userrpcinterface.QueryUserInfoResp_Data{
		Username:             result.Username,
		Email:                result.Email,
		Phone:                result.Phone,
		Profile:              result.Profile,
		LastActive:           result.LastActive,
		SignupAt:             result.SignupAt.String(),
	}

	return nil
}
