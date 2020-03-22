package vo

type AuthorizeReq struct {
	Username string `form:"username" bind:"required"`
	Password string `form:"password" bind:"required"`
}

type UnAuthorizeReq struct {
	Username string `uri:"username" bind:"required"`
}
