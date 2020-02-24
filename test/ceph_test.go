package test

import (
	"fmt"
	"go-disk/store/ceph"
	"io/ioutil"
	"testing"
)

func TestConn(t *testing.T) {

	bucket := ceph.GetCephBucket("files")

	data, err := bucket.Get("/ceph/3cc5aef31788f669d7791ae6e9cd16450ea06ae3")
	if err != nil {
		fmt.Printf("err : %v", err)
		return
	}
	ioutil.WriteFile("E:/go/go-disk/filestore/test.dat", data, 0755)
}
