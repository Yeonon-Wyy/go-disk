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
	} `yaml:"dataSource"`

	Mq struct {
		Rabbit struct {
			Url string `yaml:"url"`
			ExchangeName string `yaml:"exchangeName"`
			CephQueueName string `yaml:"cephQueueName"`
			CephErrQueueName string `yaml:"cephErrQueueName"`
			CephRoutingKey string `yaml:"cephRoutingKey"`
			CephErrRoutingKey string `yaml:"cephErrRoutingKey"`
		} `yaml:"rabbit"`
	} `yaml:"mq"`

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



