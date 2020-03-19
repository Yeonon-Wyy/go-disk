package db

import (
	"go-disk/services/file/dao"
	mydb "go-disk/services/file/db/mysql"
	"time"
)

func UpdateUserFilename(username, fileHash, filename string) bool {

	res := mydb.GetConn().
		Table(dao.UserFileDao{}.TableName()).
		Where(&dao.UserFileDao{Username:username, FileHash:fileHash}).
		Updates(map[string]interface{}{
		"file_name": filename,
		"last_update": time.Now(),
	}).Error

	return res != nil
}

func QueryUserFileMetas(username string, limit int) ([]dao.UserFileDao, bool) {
	var userFiles []dao.UserFileDao
	mydb.GetConn().Limit(limit).Where(&dao.UserFileDao{Username:username}).Find(&userFiles)
	return userFiles, true
}

func DeleteFileMeta(sha1 string, username string) bool {
	err := mydb.GetConn().
		Table(dao.UserFileDao{}.TableName()).
		Where(&dao.UserFileDao{FileHash:sha1, Username:username}).
		Update("status", 0).Error

	return err != nil
}
