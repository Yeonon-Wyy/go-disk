package db

import (
	"errors"
	"go-disk/common/constant"
	userdao "go-disk/services/user/dao"
	mysqldb "go-disk/services/user/db/mysql"
	"time"
)

func InsertUser(username, password string) bool {
	now := time.Now()
	user := userdao.UserDao{
		Username:   username,
		Password:   password,
		SignupAt:   &now,
		LastActive: &now,
		Status:     constant.UserStatusAvailable,
	}
	err := mysqldb.GetConn().Create(&user).Error
	return err == nil
}

func ExistUserByUsername(username string) bool {
	var count int
	mysqldb.GetConn().
		Table(userdao.UserDao{}.TableName()).
		Where(&userdao.UserDao{Username: username}).
		Count(&count)
	return count > 0
}

func QueryUser(username string) (*userdao.UserDao, error) {
	user := userdao.UserDao{}
	rows := mysqldb.GetConn().
		Where(&userdao.UserDao{Username: username}).
		First(&user).RowsAffected
	if rows <= 0 {
		return nil, errors.New("can't found this user")
	}

	return &user, nil
}
