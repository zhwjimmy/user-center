package asynq

import (
	"context"
	"fmt"

	"github.com/zhwjimmy/user-center/internal/infrastructure/tasks/interfaces"
	"go.uber.org/zap"
)

// asynqService Asynq 服务实现
type asynqService struct {
	client interfaces.Client
	server interfaces.Server
	config *AsynqConfig
	logger *zap.Logger
}

// NewAsynqService 创建 Asynq 服务
func NewAsynqService(cfg *AsynqConfig, handlerFactory interfaces.HandlerFactory, logger *zap.Logger) (interfaces.Service, error) {
	// 创建客户端
	client, err := NewAsynqClient(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create asynq client: %w", err)
	}

	// 创建服务器
	server, err := NewAsynqServer(cfg, logger)
	if err != nil {
		client.Close() // 清理已创建的客户端
		return nil, fmt.Errorf("failed to create asynq server: %w", err)
	}

	// 通过工厂创建并注册处理器
	handlers, err := handlerFactory.CreateHandlers(logger)
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to create handlers: %w", err)
	}

	// 注册所有处理器
	for _, handler := range handlers {
		server.RegisterHandler(handler.GetTaskType(), handler)
	}

	as := &asynqService{
		client: client,
		server: server,
		config: cfg,
		logger: logger,
	}

	logger.Info("Asynq service created successfully")
	return as, nil
}

// GetClient 获取队列客户端
func (as *asynqService) GetClient() interfaces.Client {
	return as.client
}

// GetServer 获取队列服务器
func (as *asynqService) GetServer() interfaces.Server {
	return as.server
}

// Start 启动队列服务
func (as *asynqService) Start(ctx context.Context) error {
	as.logger.Info("Starting Asynq service")

	// 启动服务器
	if err := as.server.Start(ctx); err != nil {
		return fmt.Errorf("failed to start asynq server: %w", err)
	}

	as.logger.Info("Asynq service started successfully")
	return nil
}

// Stop 停止队列服务
func (as *asynqService) Stop() error {
	as.logger.Info("Stopping Asynq service")

	var errors []error

	// 停止服务器
	if err := as.server.Stop(); err != nil {
		errors = append(errors, fmt.Errorf("failed to stop asynq server: %w", err))
	}

	// 关闭客户端
	if err := as.client.Close(); err != nil {
		errors = append(errors, fmt.Errorf("failed to close asynq client: %w", err))
	}

	as.logger.Info("Asynq service stopped")

	if len(errors) > 0 {
		return fmt.Errorf("errors during shutdown: %v", errors)
	}

	return nil
}
