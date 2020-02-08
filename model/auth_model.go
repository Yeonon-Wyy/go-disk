package model

type AuthReq struct {
	Username string `form:"username" bind:"required"`
	Token string `form:"token" bind:"required"`
}
