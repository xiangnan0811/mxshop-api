package initialize

import (
	"fmt"
	"strconv"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/xiangnan0811/mxshop-api/user-web/global"
)

func GetEnv(env string) interface{} {
	viper.AutomaticEnv()
	val := viper.Get(env)
	fmt.Println("env:", env, "val:", val)
	return val
}

func InitConfig() {
	// tencent cloud sms config from env
	global.ServerConfig.TencentSmsInfo.SecretId = GetEnv("TENCENTCLOUD_SECRET_ID").(string)
	global.ServerConfig.TencentSmsInfo.SecretKey = GetEnv("TENCENTCLOUD_SECRET_KEY").(string)

	// redis config
	global.ServerConfig.RedisInfo.Host = GetEnv("MXSHOP_REDIS_HOST").(string)
	global.ServerConfig.RedisInfo.Port, _ = strconv.Atoi(GetEnv("MXSHOP_REDIS_PORT").(string))
	global.ServerConfig.RedisInfo.Password = GetEnv("MXSHOP_REDIS_PASSWORD").(string)

	// config file
	debug := GetEnv("MXSHOP_DEBUG")
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("user-web/%s-pro.yaml", configFilePrefix)
	if debug == "true" {
		configFileName = fmt.Sprintf("user-web/%s-dev.yaml", configFilePrefix)
	}

	v := viper.New()
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := v.Unmarshal(global.ServerConfig); err != nil {
		panic(err)
	}
	zap.S().Infof("配置信息：%v", global.ServerConfig)

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		zap.S().Warnf("配置文件发生变化：%s", e.Name)
		_ = v.ReadInConfig()
		_ = v.Unmarshal(global.ServerConfig)
		zap.S().Infof("配置信息：%v", global.ServerConfig)
	})
}
