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

	Store struct {
		Ceph struct {
			AccessKey string `yaml:"accessKey"`
			SecretKey string `yaml:"secretKey"`
			RegionName string `yaml:"regionName"`
			Endpoint string `yaml:"endpoint"`
			S3LocationConstraint bool `yaml:"S3LocationConstraint"`
			S3LowercaseBucket bool `yaml:"S3LowercaseBucket"`
			FileStoreBucketName string `yaml:"fileStoreBucketName"`
			PutBinDataContentType string `yaml:"putBinDataContentType"`
			FilePathPrefix string `yaml:"filePathPrefix"`
		} `yaml:"ceph"`
	} `yaml:"store"`

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
	} `yaml:"micro"`
}



