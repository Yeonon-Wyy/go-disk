package main

import (
	"go-disk/services/apigw/router"
	"log"
)


func main() {
	err := router.Router().Run(":8081")
	if err != nil {
		log.Fatal(err)
	}
}



