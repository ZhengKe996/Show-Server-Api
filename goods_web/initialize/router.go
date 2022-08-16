package initialize

import (
	"github.com/gin-gonic/gin"
	"server_api/goods_web/middlewares"
	"server_api/goods_web/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	Router.Use(middlewares.Cors()) // 配置跨域
	ApiGroup := Router.Group("/g/v1")

	router.InitGoodsRouter(ApiGroup)

	return Router
}
