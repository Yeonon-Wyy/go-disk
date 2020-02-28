package test

import (
	"fmt"
	"go-disk/store/ceph"
	"testing"
)

func TestConn(t *testing.T) {

	bucket := ceph.GetCephBucket("files")
	res,err := bucket.List("", "", "", 10)
	if err != nil {
		fmt.Printf("err : %v", err)
		return
	}
	fmt.Println(res.Contents)
}
