package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-disk/common"
	userdb "go-disk/db"
	"go-disk/model"
	"go-disk/utils"
	"log"
	"net/http"
	"time"
)

const (
	pwdSalt = "1104459"
)

var (
	//TODO: just temp
	tokenCache = make(map[string]string)
)

type UserServiceHandler struct {
	BashPath string
}

func (u UserServiceHandler) Init(group *gin.RouterGroup) {
	group.POST("/register", userRegister())
	group.POST("/login", userLogin())
}

func userRegister() gin.HandlerFunc {
	return func(context *gin.Context) {
		var req model.UserRegisterReq
		if err := context.ShouldBind(&req); err != nil {
			log.Printf("request parameters error : %v", err)
			context.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		if userdb.ExistUserByUsername(req.Username) {
			log.Printf("user already registered : %s", req.Username)
			context.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeUserAlreadyRegistered, nil))
			return
		}

		pwd := utils.Sha1([]byte(req.Password + pwdSalt))

		if !userdb.InsertUser(req.Username, pwd) {
			context.JSON(http.StatusInternalServerError,
				common.NewServiceResp(common.RespCodeUserRegisterError, nil))
			return
		}

		context.JSON(http.StatusOK,
			common.NewServiceResp(common.RespCodeSuccess, nil))

	}
}

func userLogin() gin.HandlerFunc {
	return func(context *gin.Context) {
		var req model.UserLoginReq
		if err := context.ShouldBind(&req); err != nil {
			log.Printf("request parameters error : %v", err)
			context.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		//TODO: just temp token
		if _, ok := tokenCache[req.Username]; ok {
			log.Printf("user already login")
			context.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeUserAlreadyLogin, nil))
			return
		}

		exist := userdb.ExistUserByUsernameAndPassword(req.Username, utils.Sha1([]byte(req.Password + pwdSalt)))
		if !exist {
			context.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeUserNotFound, nil))
			return
		}

		token := genToken(req.Username)

		tokenCache[req.Username] = token

		context.JSON(http.StatusOK,
			common.NewServiceResp(common.RespCodeSuccess, token))
	}
}

func genToken(username string) string {
	ts := fmt.Sprintf("%x", time.Now().Unix())
	return utils.MD5([]byte(username + ts + "_tokensalt")) + ts[:8]
}



