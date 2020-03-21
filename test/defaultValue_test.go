package test

import (
	"fmt"
	"go-disk/common/utils"
	"reflect"
	"testing"
)

type Config struct {

	//db
	DataSource struct {
		Mysql struct {
			Username string `yaml:"username"`
			Password string `yaml:"password"`
			Host string `yaml:"host"`
			Port int `yaml:"port"`
			TimeLoc string `yaml:"timeLoc"`
			Database string `yaml:"database"`
			MaxIdle int `yaml:"maxIdle" default:"10"`
			MaxOpenConn int `yaml:"maxOpenConn" default:"100"`
			MaxLifeTime int `yaml:"maxLifeTime" default:"720"`
		} `yaml:"mysql"`

		Redis struct {
			Addr string `yaml:"addr"`
			Password string `yaml:"password"`
			Database int `yaml:"database"`
			TokenExpireTime int `yaml:"tokenExpireTime"`
		} `yaml:"redis"`
	} `yaml:"dataSource"`

	Business struct {
		UserPasswordSalt string `yaml:"userPasswordSalt"`
	} `yaml:"business"`

	Micro struct {
		Registration struct {
			Consul struct {
				Addr string `yaml:"addr"`
			} `yaml:"consul"`
		} `yaml:"registration"`
	} `yaml:"micro"`
}

type Test struct {
	A struct {
		Name string `yaml:"userPasswordSalt" default:"haha"`
	}
}

func TestDefaultValue(_ *testing.T) {
	conf := &Config{}
	ct := reflect.TypeOf(*conf)
	elements := reflect.ValueOf(conf)

	utils.SetConfigDefaultValue(ct, elements)
	fmt.Println(*conf)
}
