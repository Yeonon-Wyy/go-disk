package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"go-disk/services/auth/config"
	"log"
	"net/url"
	"os"
	"time"
)

const (
	DriverName = "mysql"
)

var (
	mysqlConfig = config.Conf.DataSource.Mysql
	mysqlDB *gorm.DB
)

func GetConn() *gorm.DB {
	if mysqlDB != nil {
		return mysqlDB
	}

	dbUrl := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=%s",
		mysqlConfig.Username,
		mysqlConfig.Password,
		mysqlConfig.Host,
		mysqlConfig.Port,
		mysqlConfig.Database,
		url.QueryEscape(mysqlConfig.TimeLoc))

	mysqlDB, err := gorm.Open(DriverName, dbUrl)
	if err != nil {
		log.Printf("connect to mysql error : %v", err)
		os.Exit(1)
	}

	if mysqlConfig.MaxIdle > 0 {
		mysqlDB.DB().SetMaxIdleConns(mysqlConfig.MaxIdle)
	}

	if mysqlConfig.MaxOpenConn > 0 {
		mysqlDB.DB().SetMaxOpenConns(mysqlConfig.MaxOpenConn)
	}

	if mysqlConfig.MaxLifeTime > 0 {
		d := time.Duration(mysqlConfig.MaxLifeTime) * time.Second
		mysqlDB.DB().SetConnMaxLifetime(d)
	}

	return mysqlDB

}