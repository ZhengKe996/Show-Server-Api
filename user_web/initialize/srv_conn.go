package initialize

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"server_api/user_web/global"
	"server_api/user_web/proto"
)

// InitSrvConn 初始化Srv连接  从注册中心获取到用户服务的消息（已经事先创立好了连接，后续不用再次进行tcp的三次握手）
func InitSrvConn() {
	// 服务注册
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)

	userSrvHost := ""
	userSrvPort := 0
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf("Service == \"%s\"", global.ServerConfig.UserSrvInfo.Name))
	if err != nil {
		panic(err)
	}
	for _, value := range data {
		userSrvHost = value.Address
		userSrvPort = value.Port
		break
	}
	if userSrvHost == "" {
		zap.S().Fatalf("[InitSrvConn] 连接 【用户服务】失败")
	}
	// 拨号连接用户 GRPC 服务
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", userSrvHost, userSrvPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 【用户服务失败】", "msg", err.Error())
	}

	// 生成grpc的client 调用接口
	userClient := proto.NewUserClient(userConn)
	global.UserClient = userClient
}
