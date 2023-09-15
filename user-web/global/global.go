package global

import (
	ut "github.com/go-playground/universal-translator"

	"github.com/xiangnan0811/mxshop-api/user-web/config"
)

var (
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
	JwtConfig    *config.JwtConfig    = &config.JwtConfig{}
	Trans        ut.Translator
)
