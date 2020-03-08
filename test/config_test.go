package test

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"testing"
)

func TestYamlConf(t *testing.T) {
	var c conf
	con := c.getConf()
	fmt.Println(c)
	fmt.Println(con.DataSource.Mysql)
}

type conf struct {
	DataSource struct{
		Mysql struct {
			Username string `yaml:"username"`
		} `yaml:"mysql"`
	} `yaml:"dataSource"`
}



func (c *conf) getConf() *conf {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Println(err)
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Println(err)
	}
	return c
}
