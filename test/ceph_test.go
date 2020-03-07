package test

import (
	"fmt"
	"go-disk/store/ceph"
	"log"
	"testing"
)

func TestConn(t *testing.T) {

	bucket := ceph.GetCephBucket("files")
	res,err := bucket.List("", "", "", 10)
	if err != nil {
		fmt.Printf("err : %v", err)
		return
	}

	fmt.Println(len(res.Contents))
}

func TestDeleteAllFile(t *testing.T) {
	bucket := ceph.GetCephBucket("files")

	res, err := bucket.List("", "", "", 10)
	if err != nil {
		fmt.Printf("err : %v", err)
		return
	}

	for _, content := range res.Contents {
		log.Println(content.Key)
		bucket.Del(content.Key)
	}

}
