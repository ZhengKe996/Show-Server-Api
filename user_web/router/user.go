package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"server_api/user_web/api"
)

// InitUserRouter 配置用户相关的url
func InitUserRouter(Router *gin.RouterGroup) {
	UserGroup := Router.Group("user")
	zap.S().Info("配置用户相关的url")
	{
		UserGroup.GET("list", api.GetUserList)
	}
}
