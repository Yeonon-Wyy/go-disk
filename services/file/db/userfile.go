package db

import (
	"go-disk/services/file/dao"
	mydb "go-disk/services/file/db/mysql"
	"log"
)

const (
	updateUserFilenameStatement = "UPDATE tbl_user_file SET file_name = ? WHERE user_name = ? AND file_sha1 = ?"

	QueryUserFileStatement = "SELECT file_sha1,file_size,file_name,upload_at,last_update FROM tbl_user_file WHERE user_name=? LIMIT ?"
)

func UpdateUserFilename(username, fileHash, filename string) bool {
	return execSql(updateUserFilenameStatement, filename, username, fileHash)
}

func QueryUserFileMetas(username string, limit int) ([]dao.UserFileDao, bool) {
	statement, err := mydb.DBConn().Prepare(QueryUserFileStatement)
	if err != nil {
		log.Printf("Failed to prepare statement : %v", err)
		return nil, false
	}

	rows, err := statement.Query(username, limit)
	if err != nil {
		log.Printf("Failed to query : %v", err)
		return nil, false
	}

	var userFiles []dao.UserFileDao
	for rows.Next() {
		userFile := dao.UserFileDao{}
		err := rows.Scan(&userFile.FileHash, &userFile.FileSize, &userFile.Filename, &userFile.UploadAt, &userFile.LastUpdate)
		if err != nil {
			log.Printf("Failed to scan row : %v", err)
			break
		}

		userFile.Username = username
		userFiles = append(userFiles, userFile)
	}

	return userFiles, true

}
