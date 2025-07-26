package interfaces

import (
	"context"
)

// Service 队列服务接口
type Service interface {
	// GetClient 获取队列客户端
	GetClient() Client
	
	// GetServer 获取队列服务器
	GetServer() Server
	
	// Start 启动队列服务
	Start(ctx context.Context) error
	
	// Stop 停止队列服务
	Stop() error
} 