package interfaces

import (
	"go.uber.org/zap"
)

// HandlerFactory 任务处理器工厂接口
type HandlerFactory interface {
	// CreateHandlers 创建所有任务处理器
	CreateHandlers(logger *zap.Logger) ([]Handler, error)
	
	// CreateHandler 创建指定类型的任务处理器
	CreateHandler(taskType string, logger *zap.Logger) (Handler, error)
} 