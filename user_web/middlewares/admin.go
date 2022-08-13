package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server_api/user_web/models"
)

// IsAdminAuth 验证用户权限
func IsAdminAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, _ := ctx.Get("claims")
		currentUser := claims.(*models.CustomClaims)

		if currentUser.AuthorityId != 0 {
			ctx.JSON(http.StatusForbidden, gin.H{"message": "无权限❎"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
