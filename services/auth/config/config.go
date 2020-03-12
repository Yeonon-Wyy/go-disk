package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

var Conf *Config

func init() {
	yamlFile, err := ioutil.ReadFile("./config/config.yaml")
	if err != nil {
		log.Println(err)
	}

	err = yaml.Unmarshal(yamlFile, &Conf)
	if err != nil {
		fmt.Println(err)
	}
}

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



