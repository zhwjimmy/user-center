package kafka

import (
	"context"

	"github.com/zhwjimmy/user-center/internal/infrastructure/messaging/interfaces"
	"go.uber.org/zap"
)

// kafkaConsumer Kafka消费者实现
type kafkaConsumer struct {
	manager *ConsumerManager
	logger  *zap.Logger
}

// NewKafkaConsumer 创建 Kafka 消费者
// 通过依赖注入接收 handlers，避免基础设施层依赖业务逻辑
func NewKafkaConsumer(cfg *KafkaClientConfig, handlers []interfaces.Handler, logger *zap.Logger) (interfaces.Consumer, error) {
	manager, err := NewConsumerManager(cfg, handlers, logger)
	if err != nil {
		return nil, err
	}

	return &kafkaConsumer{
		manager: manager,
		logger:  logger,
	}, nil
}

// Start 启动所有消费者
func (c *kafkaConsumer) Start(ctx context.Context) error {
	return c.manager.Start(ctx)
}

// Stop 停止所有消费者
func (c *kafkaConsumer) Stop() error {
	return c.manager.Stop()
}
