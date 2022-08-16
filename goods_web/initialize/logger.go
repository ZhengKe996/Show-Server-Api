package initialize

import "go.uber.org/zap"

// InitLogger 日志输出
func InitLogger() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
}
