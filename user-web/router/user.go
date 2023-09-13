package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"github.com/xiangnan0811/mxshop-api/user-web/api"
)

func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user")
	zap.S().Infow("配置用户相关的url")
	{
		UserRouter.GET("list", api.GetUserList)
	}
}
