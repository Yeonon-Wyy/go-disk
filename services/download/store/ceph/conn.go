package ceph

import (
	"go-disk/common/log4disk"
	"go-disk/services/transfer/config"
	"gopkg.in/amz.v1/aws"
	"gopkg.in/amz.v1/s3"
	"os"
)

var (
	CephConfig = config.Conf.Store.Ceph
)

var cephConn *s3.S3

func init() {
	//create bucket
	createCephBucket(CephConfig.FileStoreBucketName)
}

func Conn() *s3.S3 {
	if cephConn != nil {
		return cephConn
	}

	accessKey := CephConfig.AccessKey
	secretKey := CephConfig.SecretKey

	auth := aws.Auth{
		AccessKey:accessKey,
		SecretKey:secretKey,
	}

	region := aws.Region{
		Name: CephConfig.RegionName,
		EC2Endpoint: CephConfig.Endpoint,
		S3Endpoint: CephConfig.Endpoint,
		S3BucketEndpoint: "",
		S3LocationConstraint: CephConfig.S3LocationConstraint,
		S3LowercaseBucket: CephConfig.S3LowercaseBucket,
		Sign: aws.SignV2,
	}

	return s3.New(auth, region)
}

func GetCephBucket(bucket string) *s3.Bucket {
	return Conn().Bucket(bucket)
}

func createCephBucket(bucketName string)  {
	bucket := GetCephBucket(bucketName)
	_, err := bucket.List("","","", 100)
	if err == nil {
		log4disk.W("the bucket {%s} already exist, no need create", bucketName)
		return
	}

	err = bucket.PutBucket(s3.PublicRead)
	if err != nil {
		log4disk.E("failed to create bucket {%s} : %v", bucketName, err)
		os.Exit(1)
	}
}


