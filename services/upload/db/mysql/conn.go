package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"go-disk/services/file/config"
	"net/url"
	"os"
	"time"
)

const (
	DriverName = "mysql"
)

var (
	mysqlConfig = config.Conf.DataSource.Mysql
	mysqlDB     *gorm.DB
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
		log4disk.E("connect to mysql error : %v", err)
		os.Exit(1)
	}

	mysqlDB.DB().SetMaxOpenConns(mysqlConfig.MaxOpenConn)
	mysqlDB.DB().SetConnMaxLifetime(time.Duration(mysqlConfig.MaxLifeTime) * time.Second)
	mysqlDB.DB().SetMaxIdleConns(mysqlConfig.MaxIdle)

	return mysqlDB

}
