package db

import (
	"go-disk/services/file/dao"
	mydb "go-disk/services/file/db/mysql"
	"log"
	"time"
)

func UpdateUserFilename(username, fileHash, filename string) bool {

	err := mydb.GetConn().
		Table(dao.UserFileDao{}.TableName()).
		Where(&dao.UserFileDao{Username:username, FileHash:fileHash}).
		Updates(map[string]interface{}{
		"file_name": filename,
		"last_update": time.Now(),
	}).Error

	return err == nil
}

func QueryUserFileMetas(username string, limit int) ([]dao.UserFileDao, bool) {
	var userFiles []dao.UserFileDao
	mydb.GetConn().Limit(limit).Where(&dao.UserFileDao{Username:username}).Find(&userFiles)
	return userFiles, true
}

func DeleteFileMeta(sha1 string, filename, username string) bool {

	uf := dao.UserFileDao{}
	rowAffect := mydb.GetConn().
		Where(&dao.UserFileDao{FileHash:sha1, Username:username, FileName:filename}).
		Select("id").
		Find(&uf).RowsAffected
	if rowAffect <= 0 || uf.Id < 0{
		log.Printf("can't find this record")
		return false
	}

	err := mydb.GetConn().
		Delete(&dao.UserFileDao{Id: uint(uf.Id)})
	return err == nil
}

func ExistByFileHashAndUsername(fileHash string, username string) bool {
	var count int
	mydb.GetConn().
		Table(dao.UserFileDao{}.TableName()).
		Where(&dao.UserFileDao{Username:username, FileHash:fileHash}).
		Count(&count)
	return count > 0
}
