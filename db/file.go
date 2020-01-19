package db

import (
	"go-disk/common"
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

	updateFileStatement = "UPDATE tbl_file SET file_name = ? where file_sha1 = ?"

	queryFileStatement = "SELECT file_sha1,file_name,file_size,file_addr,create_at,update_at,status FROM tbl_file " +
		"WHERE file_sha1=? AND status=1 limit 1"

	deleteFileStatement = "UPDATE tbl_file SET status = ? where file_sha1 = ?"
)

func OnFileUploadFinished(sha1, filename, location string, size int64, status int, uploadAt, updateAt time.Time) bool {
	statement, err := mysql.DBConn().Prepare(insertFileStatement)
	if err != nil {
		log.Printf("Failed to prepare statement : %v", err)
		return false
	}

	defer statement.Close()

	result, err := statement.Exec(sha1,filename,size,location, status, uploadAt, updateAt)
	if err != nil {
		log.Printf("failed to execute statemnt : %v", err)
		return false
	}

	return validateRow(result, FileTableName)
}

func OnFileUpdateFinished(sha1, filename string) bool {
	statement, err := mysql.DBConn().Prepare(updateFileStatement)
	if err != nil {
		log.Printf("Failed to prepare statement : %v", err)
		return false
	}
	defer statement.Close()

	result, err := statement.Exec(filename, sha1)
	if err != nil {
		log.Printf("failed to execute statemnt : %v", err)
		return false
	}

	return validateRow(result, FileTableName)
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
	statement, err := mysql.DBConn().Prepare(deleteFileStatement)
	if err != nil {
		log.Printf("Failed to prepare statement : %v", err)
		return false
	}
	defer statement.Close()

	result, err := statement.Exec(common.FileStatusDelete, sha1)

	return validateRow(result, FileTableName)
}






