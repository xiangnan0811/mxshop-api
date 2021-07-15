package main

import (
	"fmt"
	"go.uber.org/zap"
	"mxshop-api/user-web/initialize"
)

func main() {
	port := 8021
	// 1. 初始化 logger
	initialize.InitLogger()
	// 2. 初始化 routers
	Router := initialize.Routers()
	// 3. 启动
	zap.S().Infof("启动服务器，端口：%d", port)
	if err := Router.Run(fmt.Sprintf(":%d", port)); err != nil {
		zap.S().Panic("启动失败：", err.Error())
	}
}
