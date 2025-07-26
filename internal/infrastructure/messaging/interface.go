package messaging

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

// Producer 生产者接口
type Producer interface {
	PublishUserEvent(ctx context.Context, event interface{}) error
	PublishUserEventAsync(ctx context.Context, event interface{}) error
	Close() error
}

// Consumer 消费者接口
type Consumer interface {
	Start(ctx context.Context) error
	Stop() error
}

// Event 通用事件接口
type Event interface {
	GetTopic() string
	GetEventType() string
	GetUserID() string
	GetRequestID() string
	GetTimestamp() string
}
