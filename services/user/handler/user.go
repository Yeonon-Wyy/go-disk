package handler

import (
	"context"
	"go-disk/common"
	"go-disk/config"
	userdb "go-disk/db"
	"go-disk/services/user/proto"
	"go-disk/utils"
)

type UserHandler struct {

}

func (u UserHandler) UserRegister(ctx context.Context, req *proto.RegisterReq, resp *proto.RegisterResp) error {


	if userdb.ExistUserByUsername(req.Username) {
		resp.Code = int64(common.RespCodeUserNotFound.Code)
		resp.Message = common.RespCodeUserNotFound.Message
		return nil
	}

	pwd := utils.Sha1([]byte(req.Password + config.PwdSalt))

	if !userdb.InsertUser(req.Username, pwd) {
		resp.Code = int64(common.RespCodeUserRegisterError.Code)
		resp.Message = common.RespCodeUserRegisterError.Message
		return nil
	}

	resp.Code = int64(common.RespCodeSuccess.Code)
	resp.Message = common.RespCodeSuccess.Message
	return nil
}
