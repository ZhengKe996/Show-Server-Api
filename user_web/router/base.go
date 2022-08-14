package router

import (
	"github.com/gin-gonic/gin"
	"server_api/user_web/api"
)

func InitBaseRouter(Router *gin.RouterGroup) {
	BaseRouter := Router.Group("base")
	{
		BaseRouter.GET("captcha", api.GetCaptcha) // 图形验证码接口
		BaseRouter.POST("send_sms", api.SendSms)  // 短信验证码接口

	}
}
