package model

type UserRegisterReq struct {
	Username string `form:"username" bind:"required"`
	Password string `form:"password" bind:"required"`
}

type UserLoginReq struct {
	Username string `form:"username" bind:"required"`
	Password string `form:"password" bind:"required"`
}

type UserQueryReq struct {
	Username string `form:"username" bind:"required"`
	Token string `form:"token" bind:"required"`
}
