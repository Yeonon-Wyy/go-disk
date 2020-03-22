package db

import (
	"go-disk/services/upload/dao"
	mydb "go-disk/services/upload/db/mysql"
	"time"
)

func OnFileUploadFinished(sha1, filename, location string, size int64, status int) bool {
	//TODO: 状态值需要做成枚举，否则难以维护
	now := time.Now()
	tf := dao.TableFileDao{
		FileHash: sha1,
		Filename: filename,
		FileSize: size,
		FileAddr: location,
		CreateAt: &now,
		UpdateAt: &now,
		Status:   status,
	}

	err := mydb.GetConn().Create(&tf).Error

	return err == nil

}

func GetStatus(sha1 string) int {
	tf := dao.TableFileDao{}
	mydb.GetConn().Where(&dao.TableFileDao{FileHash: sha1}).Select("status").First(&tf)
	return tf.Status
}

func ExistFile(sha1 string) bool {
	var count int
	mydb.GetConn().
		Table(dao.TableFileDao{}.TableName()).
		Where(&dao.TableFileDao{FileHash: sha1}).
		Count(&count)
	return count > 0
}
