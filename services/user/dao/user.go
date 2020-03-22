package dao

import "time"

type UserDao struct {
	Id             uint       `gorm:"column:id"`
	Username       string     `gorm:"column:user_name"`
	Password       string     `gorm:"column:user_pwd"`
	Email          string     `gorm:"column:email"`
	Phone          string     `gorm:"column:phone"`
	Profile        string     `gorm:"column:profile"`
	EmailValidated bool       `gorm:"column:email_validated"`
	PhoneValidated bool       `gorm:"column:phone_validated"`
	LastActive     *time.Time `gorm:"column:last_active"`
	SignupAt       *time.Time `gorm:"column:signup_at"`
	Status         int        `gorm:"column:status"`
}

func (UserDao) TableName() string {
	return "tbl_user"
}
