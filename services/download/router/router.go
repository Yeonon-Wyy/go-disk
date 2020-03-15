package router

import (
	"github.com/gin-gonic/gin"
	"go-disk/services/download/api"
	"go-disk/services/download/interceptor"
)

func Router() *gin.Engine {
	route := gin.Default()

	downloadGroup := route.Group("/files")

	downloadGroup.Use(interceptor.AuthorizeInterceptor())
	downloadGroup.GET("/download", api.DownloadHandler())

	return route
}
