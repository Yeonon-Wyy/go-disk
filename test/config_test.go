package test

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestGenCode(t *testing.T) {

	content := "package test\n\n"

	content += "type User" + " struct {\n"
	fieldJson := `json:"id"`

	content += "	" + "id" + " " + "int64" + " `" + fieldJson + "` " + "\n"

	content += "}"

	filename := "gen1.go"
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm) //打开文件
	if err != nil {
		log.Printf("open file error : %v", err)
		return
	}

	n, err := io.WriteString(f, content)
	if err != nil {
		log.Printf("%v", err)
		return
	}

	log.Printf("success : %d", n)

	f.Close()
}

func TestYamlConf(t *testing.T) {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Println(err)
	}

	fmt.Println(yamlFile)

}

type conf struct {
	DataSource struct {
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
