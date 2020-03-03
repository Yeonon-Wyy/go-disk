package router

import (
	"github.com/gin-gonic/gin"
	"go-disk/midware"
	"go-disk/services/upload/api"
)

func Router() *gin.Engine {
	router := gin.Default()

	router.Use(midware.Cors())

	//create and init group
	uploadGroup := router.Group("/files")
	uploadServiceRoute(uploadGroup)
	return router
}

func uploadServiceRoute(group *gin.RouterGroup) {
	group.Use(midware.AuthorizeInterceptor())
	group.POST("/upload", api.UploadFile())
	group.POST("/fastupload", api.TryFastUpload())

	group.POST("/mpupload/init", api.InitialMultipartUpload())
	group.POST("/mpupload/uppart", api.UploadPart())
	group.POST("/mpupload/complete", api.CompleteUpload())
}
