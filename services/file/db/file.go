package db

import (
	"go-disk/common/constant"
	"go-disk/services/file/dao"
	"go-disk/services/file/db/mysql"
	"log"
)



const (
	FileTableName = "tbl_file"

	updateFileNameStatement = "UPDATE tbl_file SET file_name = ? where file_sha1 = ?"

	queryFileStatement = "SELECT file_sha1,file_name,file_size,file_addr,create_at,update_at,status FROM tbl_file " +
		"WHERE file_sha1=? AND status=1 limit 1"

	deleteFileStatement = "UPDATE tbl_file SET status = ? where file_sha1 = ?"
)

func OnFileUpdateFinished(sha1, filename string) bool {
	return execSql(updateFileNameStatement, FileTableName, filename, sha1)
}


func GetFileMeta(sha1 string) (*dao.TableFile, error) {
	statement, err := mysql.DBConn().Prepare(queryFileStatement)
	if err != nil {
		log.Printf("Failed to prepare statement : %v", err)
		return nil, err
	}
	defer statement.Close()

	tf := dao.TableFile{}


	err = statement.QueryRow(sha1).Scan(
		&tf.FileSha1, &tf.Filename, &tf.FileSize, &tf.FileLocation, &tf.CreateAt, &tf.UpdateAt, &tf.Status)

	if err != nil {
		log.Printf("query error : %v", err)
		return nil, err
	}

	return &tf, nil
}

func DeleteFileMeta(sha1 string) bool {
	return execSql(deleteFileStatement, FileTableName, constant.FileStatusDelete, sha1)
}


