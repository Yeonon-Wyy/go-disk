package main

import (
	"github.com/gin-gonic/gin"
	"go-disk/handler"
	"go-disk/midware"
	"log"
)

func main() {
	router := gin.Default()
	router.Use(midware.Cors())

	group := router.Group("/files")
	(&handler.UploadServiceHandler{BashPath: group.BasePath()}).Init(group)

	if err := router.Run(":6000"); err != nil {
		log.Fatal(err)
	}
}