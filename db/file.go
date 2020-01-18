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
	insertStatement = "INSERT INTO tbl_file(`file_sha1`,`file_name`,`file_size`," +
		"`file_addr`,`status`,`create_at`,`update_at`) VALUES(?,?,?,?,?,?,?)"

	updateStatement = "UPDATE tbl_file SET file_name = ? where file_sha1 = ?"

	queryStatement = "SELECT file_sha1,file_name,file_size,file_addr,create_at,update_at,status FROM tbl_file " +
		"WHERE file_sha1=? AND status=1 limit 1"

	deleteStatement = "UPDATE tbl_file SET status = ? where file_sha1 = ?"
)

func OnFileUploadFinished(sha1, filename, location string, size int64, status int, uploadAt, updateAt time.Time) bool {
	statement, err := mysql.DBConn().Prepare(insertStatement)
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


	if rows, err := result.RowsAffected(); err == nil {
		if rows <= 0 {
			log.Printf("faied to insert to table %s, error is : %v", "tbl_file", err)
			return false
		}
		return true
	}
	return false
}

func OnFileUpdateFinished(sha1, filename string) bool {
	statement, err := mysql.DBConn().Prepare(updateStatement)
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

	if rows, err := result.RowsAffected(); err == nil {
		if rows <= 0 {
			log.Printf("faied to insert to table %s, error is : %v", "tbl_file", err)
			return false
		}
		return true
	}
	return false
}


func GetFileMeta(sha1 string) (*TableFile, error) {
	statement, err := mysql.DBConn().Prepare(queryStatement)
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
	statement, err := mysql.DBConn().Prepare(deleteStatement)
	if err != nil {
		log.Printf("Failed to prepare statement : %v", err)
		return false
	}
	defer statement.Close()

	result, err := statement.Exec(common.FileStatusDelete, sha1)

	if rows, err := result.RowsAffected(); err == nil {
		if rows <= 0 {
			log.Printf("faied to insert to table %s, error is : %v", "tbl_file", err)
			return false
		}
		return true
	}
	return false
}




