package model

import "time"

type User struct {
	Id             int
	Username       string
	Password       string
	Email          string
	Phone          string
	EmailValidated bool
	PhoneValidated bool
	SignupAt       time.Time
	LastActive     time.Time
	Profile string
	Status int
}

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

type UserQueryResp struct {
	Username string `json:"username"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Profile string `json:"profile"`
	LastActive string `json:"last_active"`
	SignupAt time.Time `json:"signup_at"`
}
