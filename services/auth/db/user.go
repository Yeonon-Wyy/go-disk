package db

import (
	"go-disk/services/auth/dao"
	mydb "go-disk/services/auth/db/mysql"
)

func ExistUserByUsernameAndPassword(username, password string) bool {
	user := dao.UserDao{}
	rowAffect := mydb.GetConn().
		Where(&dao.UserDao{Username:username, Password:password}).
		Select("id").
		First(&user).RowsAffected
	return rowAffect > 0
}
