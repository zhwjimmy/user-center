package interfaces

import (
	"context"
)

// Service 消息队列服务接口
type Service interface {
	GetProducer() Producer
	GetConsumer() Consumer
	Start(ctx context.Context) error
	Stop() error
}
