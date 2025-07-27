package interfaces

import (
	"context"

	"github.com/hibiken/asynq"
)

// Handler 任务处理器接口
type Handler interface {
	// ProcessTask 处理任务
	ProcessTask(ctx context.Context, task *asynq.Task) error
	
	// GetTaskType 获取任务类型
	GetTaskType() string
}

// HandlerFunc 任务处理函数类型
type HandlerFunc func(ctx context.Context, task *asynq.Task) error

// HandlerAdapter 将 HandlerFunc 适配为 asynq.HandlerFunc
func HandlerAdapter(fn HandlerFunc) asynq.HandlerFunc {
	return func(ctx context.Context, task *asynq.Task) error {
		return fn(ctx, task)
	}
} 