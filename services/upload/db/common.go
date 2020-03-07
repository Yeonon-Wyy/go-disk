package db

import (
	"database/sql"
	"go-disk/services/upload/db/mysql"
	"log"
)

func validateRow(result sql.Result, tableName string) bool {
	if rows, err := result.RowsAffected(); err == nil {
		if rows <= 0 {
			log.Printf("faied to insert to table %s, error is : %v", tableName, err)
			return false
		}
		return true
	}
	return false
}

func execSql(stemSql string, tableName string, args... interface{}) bool {
	statement, err := mysql.DBConn().Prepare(stemSql)
	if err != nil {
		log.Printf("Failed to prepare statement : %v", err)
		return false
	}
	defer statement.Close()

	result, err := statement.Exec(args...)
	if err != nil {
		log.Printf("failed to execute statemnt : %v", err)
		return false
	}

	return validateRow(result, tableName)
}

func exist(sqlStem string, args ...interface{}) bool {
	statement, err := mysql.DBConn().Prepare(sqlStem)
	if err != nil {
		log.Printf("Failed to prepare statement : %v", err)
		return false
	}
	defer statement.Close()

	var count int
	err = statement.QueryRow(args...).Scan(&count)
	if err != nil {
		log.Printf("can't found this row : %v", err)
		return false
	}

	return count > 0
}