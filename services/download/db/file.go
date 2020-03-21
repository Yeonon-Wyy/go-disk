package db

import (
	"errors"
	"go-disk/services/download/dao"
	mydb "go-disk/services/download/db/mysql"
)


func GetFileMeta(sha1 string) (*dao.TableFileDao, error) {

	tf := dao.TableFileDao{}
	rowAffects := mydb.GetConn().Where(&dao.TableFileDao{FileHash:sha1, Status: 1}).First(&tf).RowsAffected
	if rowAffects <= 0 {
		return nil, errors.New("can't find this file")
	}
	return &tf, nil
}






