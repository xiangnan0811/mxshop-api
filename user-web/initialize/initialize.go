package initialize

import (
	"go.uber.org/zap"

	"github.com/xiangnan0811/mxshop-api/user-web/global"
)

func init() {
	// 1. 初始化 logger
	InitLogger()
	// 2. 初始化 config
	InitConfig()
	// 3. 初始化翻译器
	if err := InitTransLators(global.ServerConfig.Lang); err != nil {
		zap.S().Panic("初始化翻译器失败：", err.Error())
	}
	// 4. 初始化自定义验证器
	InitValidators()
}
