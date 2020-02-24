package ceph

import (
	"go-disk/config"
	"gopkg.in/amz.v1/aws"
	"gopkg.in/amz.v1/s3"
	"log"
)

var cephConn *s3.S3

func init() {
	//create bucket
	createCephBucket(config.CephFileStoreBucketName)
}

func Conn() *s3.S3 {
	if cephConn != nil {
		return cephConn
	}

	accessKey := config.CephAccessKey
	secretKey := config.CephSecretKey

	auth := aws.Auth{
		AccessKey:accessKey,
		SecretKey:secretKey,
	}

	region := aws.Region{
		Name: config.CephRegionName,
		EC2Endpoint: config.CephEndpoint,
		S3Endpoint: config.CephEndpoint,
		S3BucketEndpoint: "",
		S3LocationConstraint: config.CephS3LocationConstraint,
		S3LowercaseBucket: config.CephS3LowercaseBucket,
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
		log.Printf("the bucket {%s} already exist, no need create", bucketName)
		return
	}

	err = bucket.PutBucket(s3.PublicRead)
	if err != nil {
		log.Printf("failed to create bucket {%s} : %v", bucketName, err)
		panic("failed to create bucket")
	}
}


