package asynq

import (
	"context"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/zhwjimmy/user-center/internal/infrastructure/tasks/interfaces"
	"go.uber.org/zap"
)

// asynqClient Asynq 客户端实现
type asynqClient struct {
	client *asynq.Client
	config *AsynqConfig
	logger *zap.Logger
}

// NewAsynqClient 创建 Asynq 客户端
func NewAsynqClient(cfg *AsynqConfig, logger *zap.Logger) (interfaces.Client, error) {
	client := asynq.NewClient(cfg.GetRedisClientOpt())

	ac := &asynqClient{
		client: client,
		config: cfg,
		logger: logger,
	}

	logger.Info("Asynq client created successfully",
		zap.String("redis_addr", cfg.Redis.Addr),
		zap.Int("redis_db", cfg.Redis.DB),
	)

	return ac, nil
}

// EnqueueTask 立即入队任务
func (ac *asynqClient) EnqueueTask(ctx context.Context, task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	info, err := ac.client.EnqueueContext(ctx, task, opts...)
	if err != nil {
		ac.logger.Error("Failed to enqueue task",
			zap.String("task_type", task.Type()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to enqueue task: %w", err)
	}

	ac.logger.Debug("Task enqueued successfully",
		zap.String("task_id", info.ID),
		zap.String("task_type", task.Type()),
		zap.String("queue", info.Queue),
	)

	return info, nil
}

// EnqueueTaskIn 延迟入队任务
func (ac *asynqClient) EnqueueTaskIn(ctx context.Context, task *asynq.Task, delay time.Duration) (*asynq.TaskInfo, error) {
	opts := []asynq.Option{asynq.ProcessIn(delay)}
	return ac.EnqueueTask(ctx, task, opts...)
}

// EnqueueTaskAt 在指定时间入队任务
func (ac *asynqClient) EnqueueTaskAt(ctx context.Context, task *asynq.Task, processAt time.Time) (*asynq.TaskInfo, error) {
	opts := []asynq.Option{asynq.ProcessAt(processAt)}
	return ac.EnqueueTask(ctx, task, opts...)
}

// Close 关闭客户端连接
func (ac *asynqClient) Close() error {
	ac.logger.Info("Closing Asynq client")
	ac.client.Close()
	return nil
}
