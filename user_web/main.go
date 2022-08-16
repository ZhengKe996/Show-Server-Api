package main

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"server_api/user_web/global"
	"server_api/user_web/initialize"
	"server_api/user_web/utils/register/consul"
	myvalidator "server_api/user_web/validator"
	"syscall"
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

	// 注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", myvalidator.ValidateMobile)
		// 验证器错误翻译
		_ = v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 非法的手机号码!", true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}

	client := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	serviceId := fmt.Sprintf("%s", uuid.NewV4())
	err := client.Register(global.ServerConfig.Host, global.ServerConfig.Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceId)
	if err != nil {
		zap.S().Panic("服务注册失败", err.Error())
	}

	go func() {
		zap.S().Debugf("启动服务器,端口:%d", global.ServerConfig.Port)
		if err := r.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
			zap.S().Panic("启动失败", err.Error())
		}
	}()

	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = client.DeRegister(serviceId); err != nil {
		zap.S().Info("注销失败:", err.Error())
	} else {
		zap.S().Info("注销成功:")
	}
}
