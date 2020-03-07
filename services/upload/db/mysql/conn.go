package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go-disk/services/upload/config"
	"log"
	"net/url"
)

const (
	DriverName = "mysql"
)

var (

	dbUrl = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&loc=%s&parseTime=true",
		config.FileDBUsername,
		config.FileDBPassword,
		config.FileDBHost,
		config.FileDBPort,
		config.FileDBName,
		url.QueryEscape(config.FileDBTimeLoc))
)

func DBConn() *sql.DB {
	db, err := sql.Open(DriverName, dbUrl)
	if err != nil {
		log.Fatalf("open mysql connection error : %v", err)
	}
	//log.Println(db)
	db.SetMaxOpenConns(1000)
	err = db.Ping()
	if err != nil {
		log.Fatalf("ping mysql connection error : %v", err)
	}
	return db
}