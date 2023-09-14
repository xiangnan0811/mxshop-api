package main

import (
	"fmt"
	"go.uber.org/zap"
	"github.com/xiangnan0811/mxshop-api/user-web/global"
	"github.com/xiangnan0811/mxshop-api/user-web/initialize"
)

func main() {
	// 1. 初始化 logger
	initialize.InitLogger()
	// 2. 初始化 config
	initialize.InitConfig()
	// 3. 初始化 routers
	Router := initialize.Routers()
	// 4. 初始化 translator
	if err := initialize.InitTrans("zh"); err != nil {
		zap.S().Panic("初始化翻译器失败：", err.Error())
	}
	// 5. 启动
	zap.S().Infof("启动服务器，端口：%d", global.ServerConfig.Port)
	if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		zap.S().Panic("启动失败：", err.Error())
	}
}
