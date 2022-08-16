package initialize

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server_api/user_web/middlewares"
	"server_api/user_web/router"
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
	ApiGroup := Router.Group("/u/v1")

	router.InitUserRouter(ApiGroup)
	router.InitBaseRouter(ApiGroup)

	return Router
}
