package config

type UserService struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type ServerConfig struct {
	Name          string      `mapstructure:"name"`
	Port          int         `mapstructure:"port"`
	UserSrvConfig UserService `mapstructure:"user_srv"`
	Lang          string      `mapstructure:"lauguage"`
}
