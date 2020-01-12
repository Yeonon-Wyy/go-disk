package main

import (
	"github.com/gin-gonic/gin"
	handler "go-disk/handler/files"
	"log"
)

func main() {
	router := gin.Default()

	filesGroup := router.Group("/files")

	DispatchHandlerGroup(filesGroup)

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func DispatchHandlerGroup(rgs ...*gin.RouterGroup) {
	for _, rg := range rgs {
		switch rg.BasePath() {
		case "/files":
			handler.FilesServiceHandler{BashPath:rg.BasePath()}.Init(rg)
		default:
			log.Printf("error handler group: %s", rg.BasePath())
		}
	}

}