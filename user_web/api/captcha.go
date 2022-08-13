package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
	"net/http"
)

var store = base64Captcha.DefaultMemStore

// GetCaptcha 获取验证码
func GetCaptcha(ctx *gin.Context) {
	driver := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)
	captcha := base64Captcha.NewCaptcha(driver, store)
	if id, b64s, err := captcha.Generate(); err != nil {
		zap.S().Errorw("生成验证码错误: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "生成验证码错误❎"})
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{"captchaId": id, "picPath": b64s})
	}
}
