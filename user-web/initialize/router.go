package initialize

import (
	"github.com/gin-gonic/gin"
	"mxshop-api/user-web/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()

	ApiRouter := Router.Group("/u/v1")
	router.InitUserRouter(ApiRouter)

	return Router
}
