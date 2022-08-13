package initialize

import (
	"github.com/gin-gonic/gin"
	"server_api/user_web/middlewares"
	"server_api/user_web/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	Router.Use(middlewares.Cors()) // 配置跨域
	ApiGroup := Router.Group("/u/v1")

	router.InitUserRouter(ApiGroup)
	router.InitBaseRouter(ApiGroup)

	return Router
}
