package main

import (
	"fmt"
	"go.uber.org/zap"
	"server_api/goods_web/global"
	"server_api/goods_web/initialize"
)

func main() {
	// 初始化 logger
	initialize.InitLogger()

	// 初始化 config
	initialize.InitConfig()

	// 初始化 router
	r := initialize.Routers()

	// 初始化 翻译器
	_ = initialize.InitTrans("zh")

	// 初始化 Srv连接
	initialize.InitSrvConn()

	//  本地开发环境 端口号固定，线上环境启动获取端口号
	//if port, err := utils.GetFreePort(); err == nil {
	//	global.ServerConfig.Port = port
	//}

	zap.S().Debugf("启动服务器,端口:%d", global.ServerConfig.Port)
	if err := r.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		zap.S().Panic("启动失败", err.Error())
	}
}
