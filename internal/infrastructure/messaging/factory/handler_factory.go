package factory

import (
	"github.com/zhwjimmy/user-center/internal/infrastructure/messaging/interfaces"
	"go.uber.org/zap"
)

// DefaultHandlerFactory 默认处理器工厂
// 这是一个临时的默认实现，实际应该由业务层提供
type DefaultHandlerFactory struct{}

// NewDefaultHandlerFactory 创建默认处理器工厂
func NewDefaultHandlerFactory() *DefaultHandlerFactory {
	return &DefaultHandlerFactory{}
}

// CreateHandlers 创建所有处理器
// 这个方法应该由业务层实现，这里只是占位符
func (f *DefaultHandlerFactory) CreateHandlers(cfg interfaces.HandlerConfig, logger *zap.Logger) ([]interfaces.Handler, error) {
	// 这里应该返回业务层实现的 Handler
	// 暂时返回空切片，避免编译错误
	logger.Warn("Using default handler factory - handlers should be provided by business layer")
	return []interfaces.Handler{}, nil
}
