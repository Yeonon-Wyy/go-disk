package db

import (
	"go-disk/common/constant"
	"go-disk/db/mysql"
	"log"
	"time"
)

type TableFile struct {
	FileSha1 string
	Filename string
	FileSize int64
	FileLocation string
	CreateAt time.Time
	UpdateAt time.Time
	Status int
}

const (
	FileTableName = "tbl_file"

	insertFileStatement = "INSERT INTO tbl_file(`file_sha1`,`file_name`,`file_size`," +
		"`file_addr`,`status`,`create_at`,`update_at`) VALUES(?,?,?,?,?,?,?)"

	updateFileNameStatement = "UPDATE tbl_file SET file_name = ? where file_sha1 = ?"

	updateFileLocationStatement = "UPDATE tbl_file SET file_addr = ? where file_sha1 = ?"

	queryFileStatement = "SELECT file_sha1,file_name,file_size,file_addr,create_at,update_at,status FROM tbl_file " +
		"WHERE file_sha1=? AND status=1 limit 1"

	deleteFileStatement = "UPDATE tbl_file SET status = ? where file_sha1 = ?"


	existFileStatement = "SELECT COUNT(1) AS count FROM tbl_file WHERE file_sha1 = ?"
)

func OnFileUploadFinished(sha1, filename, location string, size int64, status int, uploadAt, updateAt time.Time) bool {
	return execSql(insertFileStatement, FileTableName, sha1, filename, size, location, status, uploadAt, updateAt)
}

func OnFileUpdateFinished(sha1, filename string) bool {
	return execSql(updateFileNameStatement, FileTableName, filename, sha1)
}


func GetFileMeta(sha1 string) (*TableFile, error) {
	statement, err := mysql.DBConn().Prepare(queryFileStatement)
	if err != nil {
		log.Printf("Failed to prepare statement : %v", err)
		return nil, err
	}
	defer statement.Close()

	tf := TableFile{}


	err = statement.QueryRow(sha1).Scan(
		&tf.FileSha1, &tf.Filename, &tf.FileSize, &tf.FileLocation, &tf.CreateAt, &tf.UpdateAt, &tf.Status)

	if err != nil {
		log.Printf("query error : %v", err)
		return nil, err
	}

	return &tf, nil
}

func DeleteFileMeta(sha1 string) bool {
	return execSql(deleteFileStatement, FileTableName, constant.FileStatusDelete, sha1)
}

func ExistFile(sha1 string) bool {
	return exist(existFileStatement, sha1)
}

func UpdateFileLocation(sha1 string, location string) bool {
	return execSql(updateFileLocationStatement, FileTableName, location, sha1)
}






