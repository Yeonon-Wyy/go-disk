package db

import (
	"errors"
	"go-disk/common/constant"
	"go-disk/services/file/dao"
	mydb "go-disk/services/file/db/mysql"
)

func GetFileMeta(sha1 string) (*dao.TableFileDao, error) {
	tf := dao.TableFileDao{}
	rowAffects := mydb.GetConn().Where(&dao.TableFileDao{FileHash: sha1, Status: constant.FileStatusAvailable}).First(&tf).RowsAffected
	if rowAffects <= 0 {
		return nil, errors.New("can't find this file")
	}
	return &tf, nil
}
