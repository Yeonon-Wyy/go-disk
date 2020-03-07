package router

import (
	"github.com/gin-gonic/gin"
	"go-disk/services/apigw/auth"
	"go-disk/services/apigw/cors"
	"go-disk/services/apigw/handler"
)

func Router() *gin.Engine {
	router := gin.Default()

	router.Static("/static", "./static")
	router.StaticFile("/hom", "./static/view/home.html")

	router.Use(cors.Cors())

	//create and init group
	userGroup := router.Group("/users")
	uploadGroup := router.Group("/files/upload")
	downloadGroup := router.Group("/files/download")
	fileMetaGroup := router.Group("/files")

	userServiceRoute(userGroup)
	uploadServiceRoute(uploadGroup)
	downloadServiceRoute(downloadGroup)
	fileMetaServiceRoute(fileMetaGroup)
	return router
}

func fileMetaServiceRoute(group *gin.RouterGroup) {
	group.Use(auth.AuthorizeInterceptor())
	group.GET("/meta", handler.GetFileMeta())

	group.PUT("/meta", handler.UpdateFileMeta())
	group.POST("/meta", handler.GetFileList())

	group.DELETE("/delete", handler.DeleteFile())
}

func downloadServiceRoute(group *gin.RouterGroup) {
	group.Use(auth.AuthorizeInterceptor())
	group.GET("/endpoint", handler.GetDownloadServiceEndpoint())
}

func uploadServiceRoute(group *gin.RouterGroup) {
	group.Use(auth.AuthorizeInterceptor())
	group.GET("/endpoint", handler.GetUploadServiceEndpoint())
}

func userServiceRoute(group *gin.RouterGroup) {
	group.POST("/register", handler.RegisterUser())
	group.POST("/login", handler.UserLogin())

	group.Use(auth.AuthorizeInterceptor())
	group.GET("/info", handler.QueryUserInfo())
}


