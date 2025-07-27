package interfaces

import (
	"context"

	"github.com/hibiken/asynq"
)

// Server 队列服务器接口
type Server interface {
	// Start 启动服务器
	Start(ctx context.Context) error

	// Stop 停止服务器
	Stop() error

	// RegisterHandler 注册任务处理器
	RegisterHandler(taskType string, handler asynq.Handler)

	// RegisterHandlerFunc 注册任务处理函数
	RegisterHandlerFunc(taskType string, handler asynq.HandlerFunc)
}
