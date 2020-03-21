package db

import (
	"errors"
	"go-disk/services/file/dao"
	mydb "go-disk/services/file/db/mysql"
)

func OnFileUpdateFinished(sha1, filename string) bool {
	err := mydb.GetConn().Where(&dao.TableFileDao{FileHash:sha1}).Update("file_name", filename).Error
	return err == nil
}

func GetFileMeta(sha1 string) (*dao.TableFileDao, error) {
	tf := dao.TableFileDao{}
	rowAffects := mydb.GetConn().Where(&dao.TableFileDao{FileHash:sha1, Status: 1}).First(&tf).RowsAffected
	if rowAffects <= 0 {
		return nil, errors.New("can't find this file")
	}
	return &tf, nil
}


