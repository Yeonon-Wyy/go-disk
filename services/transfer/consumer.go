package main

import (
	"encoding/json"
	"go-disk/common/log4disk"
	"go-disk/common/mqproto"
	"go-disk/services/transfer/config"
	mydb "go-disk/services/transfer/db"
	"go-disk/services/transfer/mq"
	"go-disk/services/transfer/store/ceph"
	"gopkg.in/amz.v1/s3"
	"io/ioutil"
	"os"
)

var (
	mqConfig    = config.Conf.Mq
	storeConfig = config.Conf.Store
)

func processTransfer(msg []byte) bool {
	//parse msg

	msgData := mqproto.RabbitMessage{}
	err := json.Unmarshal(msg, &msgData)
	if err != nil {
		log4disk.E("json unmarshal error : %v", err)
		return false
	}
	//create fd from temp location
	fd, err := os.Open(msgData.SrcLocation)
	if err != nil {
		log4disk.E("open file error : %v", err)
		fd.Close()
		return false
	}
	//upload file to ceph
	bucket := ceph.GetCephBucket(storeConfig.Ceph.FileStoreBucketName)
	fd.Seek(0, 0)
	fileData, err := ioutil.ReadAll(fd)
	if err != nil {
		log4disk.E("read file error : %v", err)
		return false
	}
	err = bucket.Put(msgData.DstLocation, fileData, storeConfig.Ceph.PutBinDataContentType, s3.PublicReadWrite)
	if err != nil {
		log4disk.E("upload file to ceph error : %v", err)
		fd.Close()
		return false
	}
	//update file meta to file_table

	if suc := mydb.UpdateFileLocation(msgData.FileHash, msgData.DstLocation); !suc {
		log4disk.E("update file location error")
		fd.Close()
		return false
	}

	//close file before remove file
	fd.Close()

	//delete file from temp store location
	err = os.Remove(msgData.SrcLocation)
	if err != nil {
		log4disk.E("remove file error : %v", err)
	}

	return true

}

func main() {
	mq.RabbitConsume(
		mqConfig.Rabbit.CephQueueName,
		"transfer_ceph",
		processTransfer)
	//fmt.Printf("a")
}
