package db

import (
	"go-disk/common"
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
	statement, err := mysqldb.DBConn().Prepare(insertUserStatement)
	if err != nil {
		log.Printf("Failed to prepare statement : %v", err)
		return false
	}

	defer statement.Close()

	res, err := statement.Exec(username, password, time.Now(), common.UserStatusAvailable)
	if err != nil {
		log.Printf("failed to execute statemnt : %v", err)
		return false
	}

	return validateRow(res, userTableName)
}

func ExistUserByUsername(username string) bool {
	return existUser(existUserByUsernameStatement, username)
}

func ExistUserByUsernameAndPassword(username, password string) bool {
	return existUser(existUserByUsernameAndPasswordStatement, username, password)
}

func existUser(sqlStem string, args ...interface{}) bool {
	statement, err := mysqldb.DBConn().Prepare(sqlStem)
	if err != nil {
		log.Printf("Failed to prepare statement : %v", err)
		return false
	}
	defer statement.Close()

	var count int
	err = statement.QueryRow(args...).Scan(&count)
	if err != nil {
		log.Printf("can't found the user : %v", err)
		return false
	}

	return count > 0
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