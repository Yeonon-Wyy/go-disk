package handler

import (
	"github.com/gin-gonic/gin"
	"go-disk/auth"
	"go-disk/common"
	userdb "go-disk/db"
	"go-disk/model"
	"go-disk/utils"
	"log"
	"net/http"
)

const (
	pwdSalt = "1104459"
)

type UserServiceHandler struct {
	BashPath string
}

func (u UserServiceHandler) Init(group *gin.RouterGroup) {
	group.POST("/register", userRegister())
	group.POST("/login", userLogin())

	group.Use(auth.AuthorizeInterceptor())
	group.GET("/info", queryUserInfo())

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
		if auth.ExistToken(req.Username) {
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

		context.JSON(http.StatusOK,
			common.NewServiceResp(common.RespCodeSuccess, auth.GenToken(req.Username)))
	}
}

func queryUserInfo() gin.HandlerFunc {
	return func(context *gin.Context) {
		var req model.UserQueryReq
		if err := context.ShouldBind(&req); err != nil {
			log.Printf("request parameters error : %v", err)
			context.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeBindReParamError, nil))
			return
		}

		resp, err := userdb.QueryUser(req.Username)
		if err != nil {
			context.JSON(http.StatusBadRequest,
				common.NewServiceResp(common.RespCodeUserNotFound, nil))
			return
		}

		context.JSON(http.StatusOK,
			common.NewServiceResp(common.RespCodeSuccess, resp))
	}
}




