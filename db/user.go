package db

import (
	"go-disk/common/constant"
	mysqldb "go-disk/db/mysql"
	"go-disk/model"
	"log"
	"time"
)

const (
	userTableName = "tbl_user"

	insertUserStatement = "INSERT INTO tbl_user(`user_name`,`user_pwd`, `signup_at`, `status`) " +
		"VALUES(?,?,?,?)"

	existUserByUsernameAndPasswordStatement = "SELECT COUNT(1) AS count FROM tbl_user WHERE user_name = ? AND user_pwd = ?"

	existUserByUsernameStatement = "SELECT COUNT(1) AS count FROM tbl_user WHERE user_name = ?"

	queryBriefUserInfoStatement = "SELECT signup_at FROM tbl_user WHERE user_name = ?"
)

func InsertUser(username, password string) bool {
	return execSql(insertUserStatement, userTableName, username, password, time.Now(), constant.UserStatusAvailable)
}

func ExistUserByUsername(username string) bool {
	return exist(existUserByUsernameStatement, username)
}

func ExistUserByUsernameAndPassword(username, password string) bool {
	return exist(existUserByUsernameAndPasswordStatement, username, password)
}



func QueryUser(username string) (*model.UserQueryResp, error){
	statement, err := mysqldb.DBConn().Prepare(queryBriefUserInfoStatement)
	if err != nil {
		log.Printf("Failed to prepare statement : %v", err)
		return nil, err
	}
	defer statement.Close()

	var resp model.UserQueryResp
	statement.QueryRow(username).Scan(&resp.SignupAt)
	resp.Username = username

	return &resp, nil
}