package db

const (

	tableName = "tbl_file"

	updateFileLocationStatement = "UPDATE tbl_file SET file_addr = ? WHERE file_sha1 = ?"

)

func UpdateFileLocation(fileHash, location string) bool {
	return execSql(updateFileLocationStatement, tableName, location, fileHash)
}

