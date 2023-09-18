package config

type UserService struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type ServerConfig struct {
	Name           string                `mapstructure:"name"`
	Port           int                   `mapstructure:"port"`
	UserSrvConfig  UserService           `mapstructure:"user_srv"`
	Lang           string                `mapstructure:"lauguage"`
	JWTInfo        JwtConfig             `mapstructure:"jwt"`
	TencentSmsInfo TencentCloudSmsConfig `mapstructure:"tencent"`
	RedisInfo      RedisConfig           `mapstructure:"redis"`
	SmsInfo        SmsConfig             `mapstructure:"sms"`
}

type JwtConfig struct {
	SigningKey string `mapstructure:"key"`
	Expires    int64  `mapstructure:"expires"`
	Issuer     string `mapstructure:"issuer"`
}

type TencentCloudSmsConfig struct {
	SecretId    string
	SecretKey   string
	SmsSdkAppId string `mapstructure:"sms_sdk_appId"`
	TemplateId  string `mapstructure:"template_id"`
	SignName    string `mapstructure:"sign_name"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
}

type SmsConfig struct {
	Length     int `mapstructure:"length"`
	Expires    int `mapstructure:"expires"`
	Interval   int `mapstructure:"interval"`
	SmsRedisDB int `mapstructure:"sms_db"`
}

type CaptchaConfig struct {
	Length int `mapstructure:"length"`
	Width  int `mapstructure:"width"`
	Height int `mapstructure:"height"`
}
