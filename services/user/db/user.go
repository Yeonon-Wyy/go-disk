package db

import (
	"go-disk/common/constant"
	userdao "go-disk/services/user/dao"
	mysqldb "go-disk/services/user/db/mysql"
	"log"
	"time"
)

const (
	userTableName = "tbl_user"

	insertUserStatement = "INSERT INTO tbl_user(`user_name`,`user_pwd`, `signup_at`, `status`) " +
		"VALUES(?,?,?,?)"

	existUserByUsernameStatement = "SELECT COUNT(1) AS count FROM tbl_user WHERE user_name = ?"

	queryBriefUserInfoStatement = "SELECT signup_at FROM tbl_user WHERE user_name = ?"
)

func InsertUser(username, password string) bool {
	return execSql(insertUserStatement, userTableName, username, password, time.Now(), constant.UserStatusAvailable)
}

func ExistUserByUsername(username string) bool {
	return exist(existUserByUsernameStatement, username)
}


func QueryUser(username string) (*userdao.UserQueryDao, error){
	statement, err := mysqldb.DBConn().Prepare(queryBriefUserInfoStatement)
	if err != nil {
		log.Printf("Failed to prepare statement : %v", err)
		return nil, err
	}
	defer statement.Close()

	var resp userdao.UserQueryDao
	statement.QueryRow(username).Scan(&resp.SignupAt)
	resp.Username = username

	return &resp, nil
}