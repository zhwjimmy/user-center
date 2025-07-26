package interfaces

import (
	"go.uber.org/zap"
)

// HandlerFactory 处理器工厂接口
// 由业务层实现，基础设施层通过接口获取 Handler
type HandlerFactory interface {
	CreateHandlers(cfg HandlerConfig, logger *zap.Logger) ([]Handler, error)
}

// HandlerConfig 处理器配置接口
// 只包含 Handler 创建所需的最小配置信息
type HandlerConfig interface {
	GetGroupID() string
	GetTopicName(key string) string
}
