package db

import (
	"go-disk/services/upload/dao"
	mydb "go-disk/services/upload/db/mysql"
	"time"
)

func InsertUserFile(username, fileHash, filename, filePath string, fileSize int64, status int) bool {
	now := time.Now()
	uf := dao.UserFileDao{
		Username:   username,
		FileHash:   fileHash,
		FileSize:   fileSize,
		FileName:   filename,
		FilePath:   filePath,
		UploadAt:   &now,
		LastUpdate: &now,
		Status:     status,
	}
	err := mydb.GetConn().Create(&uf).Error

	return err == nil
}


