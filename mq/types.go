package mq

import "go-disk/common"

type RabbitMessage struct {
	FileHash string
	SrcLocation string
	DstLocation string
	DstStoreType common.StoreType
}
