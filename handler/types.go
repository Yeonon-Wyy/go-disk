package handler

import "github.com/gin-gonic/gin"

type ServiceHandler interface {
	 Init(group *gin.RouterGroup)
}
