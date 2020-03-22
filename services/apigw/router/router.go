package router

import (
	"github.com/gin-gonic/gin"
	"go-disk/services/apigw/interceptor"
	"go-disk/services/apigw/cors"
	"go-disk/services/apigw/api"
)

func Router() *gin.Engine {
	router := gin.Default()

	router.Static("/static", "../../static")
	router.StaticFile("/hom", "../../static/view/home.html")

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
	group.Use(interceptor.AuthorizeInterceptor())
	group.GET("/meta/:username/:file_hash", api.GetFileMeta())

	group.PUT("/meta/:username/:file_hash", api.UpdateFileMeta())
	group.GET("/meta/:username", api.GetFileList())

	group.DELETE("/meta/:username/:file_hash", api.DeleteFile())
}

func downloadServiceRoute(group *gin.RouterGroup) {
	group.Use(interceptor.AuthorizeInterceptor())
	group.GET("/endpoint", api.GetDownloadServiceEndpoint())
}

func uploadServiceRoute(group *gin.RouterGroup) {
	group.Use(interceptor.AuthorizeInterceptor())
	group.GET("/endpoint", api.GetUploadServiceEndpoint())
}

func userServiceRoute(group *gin.RouterGroup) {
	group.POST("/register", api.RegisterUser())
	group.POST("/login", api.UserLogin())

	group.Use(interceptor.AuthorizeInterceptor())
	group.GET("/info", api.QueryUserInfo())
}


