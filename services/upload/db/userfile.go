package db

import (
	"time"
)

const (
	userFileTableName = "tbl_user_file"

	insertUserFileStatement = "INSERT INTO tbl_user_file(`user_name`,`file_sha1`,`file_name`,`file_size`, `upload_at`)" +
		"VALUES(?,?,?,?,?)"
)

func InsertUserFile(username, fileHash, filename string, fileSize int64) bool {
	return execSql(insertUserFileStatement, userFileTableName, username, fileHash, filename, fileSize, time.Now())
}


