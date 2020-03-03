package main

import (
	"encoding/json"
	"go-disk/config"
	mydb "go-disk/db"
	"go-disk/mq"
	"go-disk/store/ceph"
	"gopkg.in/amz.v1/s3"
	"io/ioutil"
	"log"
	"os"
)

func processTransfer(msg []byte) bool {
	//parse msg

	msgData := mq.RabbitMessage{}
	err := json.Unmarshal(msg, &msgData)
	if err != nil {
		log.Printf("json unmarshal error : %v", err)
		return false
	}
	//create fd from temp location
	fd, err := os.Open(msgData.SrcLocation)
	if err != nil {
		log.Printf("open file error : %v", err)
		return false
	}
	//upload file to ceph
	bucket := ceph.GetCephBucket(config.CephFileStoreBucketName)
	fd.Seek(0, 0)
	fileData, err := ioutil.ReadAll(fd)
	if err != nil {
		log.Printf("read file error : %v", err)
		return false
	}
	err = bucket.Put(msgData.DstLocation, fileData, config.CephPutBinDataContentType, s3.PublicReadWrite)
	if err != nil {
		log.Printf("upload file to ceph error : %v", err)
		return false
	}
	//update file meta to file_table

	if suc := mydb.UpdateFileLocation(msgData.FileHash, msgData.DstLocation); !suc {
		log.Printf("update file location error")
		return false
	}


	return true

}

func main() {
	mq.RabbitConsume(
		config.RabbitCephQueueName,
		"transfer_ceph",
		processTransfer)
	//fmt.Printf("a")
}
