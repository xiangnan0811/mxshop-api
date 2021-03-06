package main

import (
	"fmt"
	"go.uber.org/zap"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/initialize"
)

func main() {
	// 1. 初始化 logger
	initialize.InitLogger()
	// 2. 初始化 config
	initialize.InitConfig()
	// 2. 初始化 routers
	Router := initialize.Routers()
	// 3. 启动
	zap.S().Infof("启动服务器，端口：%d", global.ServerConfig.Port)
	if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		zap.S().Panic("启动失败：", err.Error())
	}
}
