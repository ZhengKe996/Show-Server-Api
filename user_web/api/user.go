package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"net/http"
	"server_api/user_web/global"
	"server_api/user_web/global/reponse"
	"server_api/user_web/proto"
	"time"
)

// HandleGrpcError2Http 将GRPC的Code转换成Http的状态码
func HandleGrpcError2Http(err error, c *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{"message": e.Message()})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{"message": "内部错误❎"})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{"message": "参数错误❎"})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{"message": "用户服务不可用"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"message": "其他错误❎"})
			}
			return
		}
	}
}

// GetUserList 获取用户列表
func GetUserList(ctx *gin.Context) {
	// 拨号连接用户 GRPC 服务
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host, global.ServerConfig.UserSrvInfo.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 【用户服务失败】", "msg", err.Error())
	}

	// 生成grpc的client 调用接口
	userClient := proto.NewUserClient(userConn)
	response, err := userClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    0,
		PSize: 10,
	})
	if err != nil {
		zap.S().Error("[GetUserList] 【查询用户列表失败】")
		HandleGrpcError2Http(err, ctx)
		return
	}
	result := make([]interface{}, 0)
	for _, value := range response.Data {
		userResponse := reponse.UserResponse{
			ID:       value.Id,
			NikeName: value.NickName,
			Mobile:   value.Mobile,
			Birthday: reponse.JsonTime(time.Unix(int64(value.BirthDay), 0)),
			Gender:   value.Gender,
		}

		result = append(result, userResponse)
	}

	ctx.JSON(http.StatusOK, result)
}
