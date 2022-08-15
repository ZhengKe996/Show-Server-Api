package global

import (
	ut "github.com/go-playground/universal-translator"
	"server_api/user_web/config"
	"server_api/user_web/proto"
)

var (
	Trans        ut.Translator
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
	NaCosConfig  *config.NaCosConfig  = &config.NaCosConfig{}
	UserClient   proto.UserClient
)
