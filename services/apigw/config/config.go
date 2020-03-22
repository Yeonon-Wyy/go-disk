package config

import (
	"fmt"
	"go-disk/common/log4disk"
	"go-disk/common/utils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"reflect"
)

var Conf *Config

func init() {
	Conf = &Config{}
	utils.SetConfigDefaultValue(reflect.TypeOf(*Conf), reflect.ValueOf(Conf))

	yamlFile, err := ioutil.ReadFile("./config/config.yaml")
	if err != nil {
		log4disk.E("read file error : %v", err)
	}

	err = yaml.Unmarshal(yamlFile, &Conf)
	if err != nil {
		fmt.Println(err)
	}


}

type Config struct {

	Business struct {
		UploadServiceEndpoint string `yaml:"uploadServiceEndpoint"`
		FileStorePath string `yaml:"fileStorePath"`
	} `yaml:"business"`

	Micro struct {
		Registration struct {
			Consul struct {
				Addr string `yaml:"addr"`
			} `yaml:"consul"`
		} `yaml:"registration"`

		Client struct {
			Auth struct {
				ServiceName string `yaml:"serviceName"`
			} `yaml:"auth"`

			Download struct {
				ServiceName string `yaml:"serviceName"`
			} `yaml:"download"`

			File struct {
				ServiceName string `yaml:"serviceName"`
			} `yaml:"file"`

			Upload struct {
				ServiceName string `yaml:"serviceName"`
			} `yaml:"upload"`

			User struct {
				ServiceName string `yaml:"serviceName"`
			} `yaml:"user"`

		} `yaml:"client"`
	} `yaml:"micro"`
}



