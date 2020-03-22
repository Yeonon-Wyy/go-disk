package main

import (
	"go-disk/common/log4disk"
	"go-disk/services/apigw/router"
)

func main() {
	err := router.Router().Run(":8081")
	if err != nil {
		log4disk.E("start service error : %v", err)
	}
}
