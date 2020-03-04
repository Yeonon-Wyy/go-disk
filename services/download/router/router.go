package router

import (
	"github.com/gin-gonic/gin"
	"go-disk/midware"
	"go-disk/services/download/api"
)

func Router() *gin.Engine {
	route := gin.Default()

	downloadGroup := route.Group("/files")

	downloadGroup.Use(midware.AuthorizeInterceptor())
	downloadGroup.POST("/download", api.DownloadHandler())

	return route
}
