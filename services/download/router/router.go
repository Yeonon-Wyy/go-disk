package router

import (
	"github.com/gin-gonic/gin"
	"go-disk/services/download/api"
)

func Router() *gin.Engine {
	route := gin.Default()

	downloadGroup := route.Group("/files")

	downloadGroup.GET("/download", api.DownloadHandler())

	return route
}
