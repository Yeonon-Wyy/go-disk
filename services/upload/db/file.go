package db

import (
	"time"
)

const (
	FileTableName = "tbl_file"

	insertFileStatement = "INSERT INTO tbl_file(`file_sha1`,`file_name`,`file_size`," +
		"`file_addr`,`status`,`create_at`,`update_at`) VALUES(?,?,?,?,?,?,?)"

	existFileStatement = "SELECT COUNT(1) AS count FROM tbl_file WHERE file_sha1 = ?"
)

func OnFileUploadFinished(sha1, filename, location string, size int64, status int, uploadAt, updateAt time.Time) bool {
	return execSql(insertFileStatement, FileTableName, sha1, filename, size, location, status, uploadAt, updateAt)
}

func ExistFile(sha1 string) bool {
	return exist(existFileStatement, sha1)
}