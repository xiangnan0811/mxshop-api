package initialize

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/xiangnan0811/mxshop-api/user-web/middlewares"
	"github.com/xiangnan0811/mxshop-api/user-web/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()

	// 跨域问题
	Router.Use(middlewares.Cors())

	// 健康检查
	Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "ok",
		})
	})
	ApiRouter := Router.Group("/u/v1")
	router.InitUserRouter(ApiRouter) // 注册用户路由
	router.InitBaseRouter(ApiRouter) // 注册基础路由

	return Router
}
