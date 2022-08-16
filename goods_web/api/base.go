package api

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"server_api/goods_web/global"
	"strings"
)

func removeTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fileds {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

// HandleGrpcError2Http 将GRPC的Code转换成Http的状态码
func HandleGrpcError2Http(err error, c *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{"message": e.Message() + "❎"})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{"message": "内部错误❎"})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{"message": "参数错误❎"})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{"message": "商品服务不可用❎"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"message": "其他错误❎"})
			}
			return
		}
	}
}

// HandleValidatorError 处理表单验证的错误信息
func HandleValidatorError(ctx *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	ctx.JSON(http.StatusBadRequest, gin.H{
		"error": removeTopStruct(errs.Translate(global.Trans)),
	})
}
