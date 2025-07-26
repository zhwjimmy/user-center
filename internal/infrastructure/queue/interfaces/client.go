package interfaces

import (
	"context"
	"time"

	"github.com/hibiken/asynq"
)

// Client 队列客户端接口
type Client interface {
	// EnqueueTask 立即入队任务
	EnqueueTask(ctx context.Context, task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error)

	// EnqueueTaskIn 延迟入队任务
	EnqueueTaskIn(ctx context.Context, task *asynq.Task, delay time.Duration) (*asynq.TaskInfo, error)

	// EnqueueTaskAt 在指定时间入队任务
	EnqueueTaskAt(ctx context.Context, task *asynq.Task, processAt time.Time) (*asynq.TaskInfo, error)

	// Close 关闭客户端连接
	Close() error
}
