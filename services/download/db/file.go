package db

import (
	"go-disk/services/download/db/mysql"
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
	queryFileStatement = "SELECT file_sha1,file_name,file_size,file_addr,create_at,update_at,status FROM tbl_file " +
		"WHERE file_sha1=? AND status=1 limit 1"
)

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






