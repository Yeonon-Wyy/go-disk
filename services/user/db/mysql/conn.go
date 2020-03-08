package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go-disk/services/user/config"
	"log"
	"net/url"
)

const (
	DriverName = "mysql"
)

var (

	DSConfig = config.Conf.DataSource

	dbUrl = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&loc=%s&parseTime=true",
		DSConfig.Mysql.Username,
		DSConfig.Mysql.Password,
		DSConfig.Mysql.Host,
		DSConfig.Mysql.Port,
		DSConfig.Mysql.Database,
		url.QueryEscape(DSConfig.Mysql.TimeLoc))
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