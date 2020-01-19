package db

import (
	"database/sql"
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