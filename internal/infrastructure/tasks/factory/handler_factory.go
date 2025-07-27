package factory

import (
	"fmt"

	"github.com/zhwjimmy/user-center/internal/infrastructure/tasks/interfaces"
	"go.uber.org/zap"
)

// DefaultHandlerFactory 默认处理器工厂
type DefaultHandlerFactory struct {
	// 可以在这里添加依赖注入
}

// NewDefaultHandlerFactory 创建默认处理器工厂
func NewDefaultHandlerFactory() interfaces.HandlerFactory {
	return &DefaultHandlerFactory{}
}

// CreateHandlers 创建所有任务处理器
func (f *DefaultHandlerFactory) CreateHandlers(logger *zap.Logger) ([]interfaces.Handler, error) {
	// 这里返回空的处理器列表，实际的处理器将在业务层创建
	// 这样可以避免基础设施层依赖业务逻辑
	logger.Info("Creating empty handler list - handlers should be registered in business layer")
	return []interfaces.Handler{}, nil
}

// CreateHandler 创建指定类型的任务处理器
func (f *DefaultHandlerFactory) CreateHandler(taskType string, logger *zap.Logger) (interfaces.Handler, error) {
	// 这里返回 nil，实际的处理器创建逻辑在业务层
	logger.Warn("Handler creation not implemented in infrastructure layer",
		zap.String("task_type", taskType),
	)
	return nil, fmt.Errorf("handler creation not implemented for task type: %s", taskType)
}
