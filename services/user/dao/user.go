package dao

import "time"

type UserQueryDao struct {
	Username string
	Email string
	Phone string
	Profile string
	LastActive string
	SignupAt time.Time
}

