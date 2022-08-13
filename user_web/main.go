package main

import (
	"fmt"
	"go.uber.org/zap"
	"server_api/user_web/global"
	"server_api/user_web/initialize"
)

func main() {
	// 初始化 logger
	initialize.InitLogger()

	// 初始化 config
	initialize.InitConfig()

	// 初始化 router
	r := initialize.Routers()

	zap.S().Debugf("启动服务器,端口:%d", global.ServerConfig.Port)
	if err := r.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		zap.S().Panic("启动失败", err.Error())
	}
}
