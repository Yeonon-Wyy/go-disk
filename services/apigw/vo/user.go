package vo

type UserRegisterReq struct {
	Username string `form:"username" bind:"required"`
	Password string `form:"password" bind:"required"`
}

type UserQueryReq struct {
	Username string `uri:"username" bind:"required"`
}
