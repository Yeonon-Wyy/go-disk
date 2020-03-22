package db

import (
	"go-disk/services/auth/dao"
	mydb "go-disk/services/auth/db/mysql"
)

func ExistUserByUsernameAndPassword(username, password string) bool {
	var count int
	mydb.GetConn().
		Table(dao.UserDao{}.TableName()).
		Where(&dao.UserDao{Username: username, Password: password}).
		Count(&count)
	return count > 0
}
