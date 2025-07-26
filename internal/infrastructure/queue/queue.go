package queue

import (
	"github.com/zhwjimmy/user-center/internal/config"
	"github.com/zhwjimmy/user-center/internal/infrastructure/queue/asynq"
	"github.com/zhwjimmy/user-center/internal/infrastructure/queue/factory"
	"github.com/zhwjimmy/user-center/internal/infrastructure/queue/interfaces"
	"go.uber.org/zap"
)

// Service 导出队列服务接口类型
type Service = interfaces.Service

// NewAsynqClientConfig 创建 Asynq 客户端配置
func NewAsynqClientConfig(cfg *config.Config) *asynq.AsynqConfig {
	return asynq.NewAsynqConfig(cfg)
}

// NewAsynqService 创建 Asynq 服务
func NewAsynqService(cfg *config.Config, logger *zap.Logger) (Service, error) {
	asynqConfig := NewAsynqClientConfig(cfg)
	handlerFactory := factory.NewDefaultHandlerFactory()

	return asynq.NewAsynqService(asynqConfig, handlerFactory, logger)
}
