package router

import (
	"github.com/gin-gonic/gin"
	"go-disk/services/apigw/api"
	"go-disk/services/apigw/cors"
	"go-disk/services/apigw/interceptor"
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
	authGroup := router.Group("/token")

	userServiceRoute(userGroup)
	uploadServiceRoute(uploadGroup)
	downloadServiceRoute(downloadGroup)
	fileMetaServiceRoute(fileMetaGroup)
	authServiceRoute(authGroup)
	return router
}

func authServiceRoute(group *gin.RouterGroup) {
	group.POST("/", api.Authorize())

	group.Use(interceptor.AuthorizeInterceptor())
	group.DELETE("/:username", api.UnAuthorize())
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
	group.POST("/", api.RegisterUser())

	group.Use(interceptor.AuthorizeInterceptor())
	group.GET("/:username", api.QueryUserInfo())
}
