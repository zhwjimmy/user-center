package asynq

import (
	"context"
	"fmt"
	"sync"

	"github.com/hibiken/asynq"
	"github.com/zhwjimmy/user-center/internal/infrastructure/tasks/interfaces"
	"go.uber.org/zap"
)

// asynqServer Asynq 服务器实现
type asynqServer struct {
	server *asynq.Server
	mux    *asynq.ServeMux
	config *AsynqConfig
	logger *zap.Logger
	mu     sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc
}

// NewAsynqServer 创建 Asynq 服务器
func NewAsynqServer(cfg *AsynqConfig, logger *zap.Logger) (interfaces.Server, error) {
	server := asynq.NewServer(cfg.GetRedisClientOpt(), cfg.GetServerConfig())
	mux := asynq.NewServeMux()

	as := &asynqServer{
		server: server,
		mux:    mux,
		config: cfg,
		logger: logger,
	}

	logger.Info("Asynq server created successfully",
		zap.String("redis_addr", cfg.Redis.Addr),
		zap.Int("redis_db", cfg.Redis.DB),
		zap.Int("workers", cfg.Workers),
		zap.Any("queues", cfg.Queues),
	)

	return as, nil
}

// Start 启动服务器
func (as *asynqServer) Start(ctx context.Context) error {
	as.mu.Lock()
	defer as.mu.Unlock()

	if as.ctx != nil {
		return fmt.Errorf("server is already running")
	}

	as.ctx, as.cancel = context.WithCancel(ctx)

	as.logger.Info("Starting Asynq server")

	// 在 goroutine 中启动服务器
	go func() {
		if err := as.server.Run(as.mux); err != nil {
			as.logger.Error("Asynq server stopped with error", zap.Error(err))
		}
	}()

	as.logger.Info("Asynq server started successfully")
	return nil
}

// Stop 停止服务器
func (as *asynqServer) Stop() error {
	as.mu.Lock()
	defer as.mu.Unlock()

	if as.ctx == nil {
		return nil
	}

	as.logger.Info("Stopping Asynq server")

	if as.cancel != nil {
		as.cancel()
	}

	as.server.Shutdown()
	as.ctx = nil
	as.cancel = nil

	as.logger.Info("Asynq server stopped successfully")
	return nil
}

// RegisterHandler 注册任务处理器
func (as *asynqServer) RegisterHandler(taskType string, handler asynq.Handler) {
	as.mu.Lock()
	defer as.mu.Unlock()

	as.mux.Handle(taskType, handler)
	as.logger.Info("Registered task handler",
		zap.String("task_type", taskType),
	)
}

// RegisterHandlerFunc 注册任务处理函数
func (as *asynqServer) RegisterHandlerFunc(taskType string, handler asynq.HandlerFunc) {
	as.mu.Lock()
	defer as.mu.Unlock()

	as.mux.HandleFunc(taskType, handler)
	as.logger.Info("Registered task handler function",
		zap.String("task_type", taskType),
	)
}
