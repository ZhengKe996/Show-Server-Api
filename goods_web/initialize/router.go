package initialize

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server_api/goods_web/middlewares"
	"server_api/goods_web/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})
	Router.Use(middlewares.Cors()) // 配置跨域
	ApiGroup := Router.Group("/g/v1")

	router.InitGoodsRouter(ApiGroup)
	return Router
}
