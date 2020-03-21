package db

import (
	"go-disk/services/transfer/dao"
	mydb "go-disk/services/transfer/db/mysql"
)



func UpdateFileLocation(fileHash, location string) bool {
	rowAffect := mydb.GetConn().Table(dao.TableFileDao{}.TableName()).
		Where(&dao.TableFileDao{FileHash:fileHash}).
		Update(map[string]interface{}{"file_addr" : location}).RowsAffected
	return rowAffect > 0
}

