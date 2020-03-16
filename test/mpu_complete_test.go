package test

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"testing"
)

func TestFileMerge(t *testing.T) {
	fii, err := os.OpenFile("E:/444.zip", os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
 		log.Printf("open file error : %v", err)
		return
	}

	for i := 1; i <= 10; i++ {
		fileName := "E:/go/go-disk/filestore/yeonon15fcccfe6fd0de30/" + strconv.Itoa(i)
		f, err := os.OpenFile(fileName, os.O_RDONLY, os.ModePerm)
		if err != nil {
			log.Printf("open file error : %v", err)
			return
		}

		b, err := ioutil.ReadAll(f)

		fii.Write(b)
		f.Close()
	}

	fii.Close()
}
