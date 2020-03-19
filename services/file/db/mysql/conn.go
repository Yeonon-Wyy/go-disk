package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"go-disk/services/file/config"
	"log"
	"net/url"
	"os"
)

const (
	DriverName = "mysql"
)

var (
	DSConfig = config.Conf.DataSource
	mysqlDB *gorm.DB
)

func GetConn() *gorm.DB {
	if mysqlDB != nil {
		return mysqlDB
	}

	dbUrl := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=%s",
		DSConfig.Mysql.Username,
		DSConfig.Mysql.Password,
		DSConfig.Mysql.Host,
		DSConfig.Mysql.Port,
		DSConfig.Mysql.Database,
		url.QueryEscape(DSConfig.Mysql.TimeLoc))

	mysqlDB, err := gorm.Open(DriverName, dbUrl)
	if err != nil {
		log.Printf("connect to mysql error : %v", err)
		os.Exit(1)
	}

	return mysqlDB

}