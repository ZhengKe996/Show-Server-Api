package initialize

import (
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"server_api/goods_web/global"
	"server_api/goods_web/proto"
)

// InitSrvConn 初始化Srv连接  从注册中心获取到用户服务的消息（已经事先创立好了连接，后续不用再次进行tcp的三次握手）
func InitSrvConn() {
	Conn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port, global.ServerConfig.GoodsSrvInfo.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务】失败")
	}

	// 生成grpc的client 调用接口
	Client := proto.NewGoodsClient(Conn)
	global.GoodsClient = Client
}
