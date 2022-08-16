package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"server_api/goods_web/api/goods"
)

func InitGoodsRouter(Router *gin.RouterGroup) {
	GoodsGroup := Router.Group("goods")
	zap.S().Info("配置用户相关的url")
	{
		GoodsGroup.GET("", goods.List)
		GoodsGroup.POST("", goods.New)
	}
}
