package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/xiangnan0811/mxshop-api/user-web/api"
	"github.com/xiangnan0811/mxshop-api/user-web/middlewares"
)

func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user")
	zap.S().Infow("配置用户相关的url")
	{
		UserRouter.GET("list", middlewares.JWTAuth(), api.GetUserList) // 用户列表页
		UserRouter.POST("login", api.PassWordLogin)                    // 用户登录
	}
}
