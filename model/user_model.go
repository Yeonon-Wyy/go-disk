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
