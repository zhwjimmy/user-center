package kafka

import (
	"context"
	"sync"

	"github.com/zhwjimmy/user-center/internal/infrastructure/messaging/interfaces"
	"go.uber.org/zap"
)

// ConsumerManager 消费者管理器
type ConsumerManager struct {
	consumers map[string]*TopicConsumer
	logger    *zap.Logger
	wg        sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewConsumerManager 创建消费者管理器
func NewConsumerManager(cfg *KafkaClientConfig, handlers []interfaces.Handler, logger *zap.Logger) (*ConsumerManager, error) {
	consumers := make(map[string]*TopicConsumer)

	for _, handler := range handlers {
		consumer, err := NewTopicConsumer(cfg, handler, logger)
		if err != nil {
			return nil, err
		}
		consumers[handler.GetTopicName()] = consumer
	}

	return &ConsumerManager{
		consumers: consumers,
		logger:    logger,
	}, nil
}

// Start 启动所有消费者
func (cm *ConsumerManager) Start(ctx context.Context) error {
	cm.ctx, cm.cancel = context.WithCancel(ctx)

	for topic, consumer := range cm.consumers {
		topic := topic
		consumer := consumer

		cm.wg.Add(1)
		go func() {
			defer cm.wg.Done()
			if err := consumer.Start(cm.ctx); err != nil {
				cm.logger.Error("Failed to start topic consumer",
					zap.String("topic", topic),
					zap.Error(err))
			}
		}()
	}

	cm.logger.Info("Consumer manager started",
		zap.Int("topic_count", len(cm.consumers)))
	return nil
}

// Stop 停止所有消费者
func (cm *ConsumerManager) Stop() error {
	cm.logger.Info("Stopping consumer manager")
	if cm.cancel != nil {
		cm.cancel()
	}

	cm.wg.Wait()

	// 停止所有消费者
	for topic, consumer := range cm.consumers {
		if err := consumer.Stop(); err != nil {
			cm.logger.Error("Failed to stop topic consumer",
				zap.String("topic", topic),
				zap.Error(err))
		}
	}

	return nil
}

// GetConsumer 获取指定 Topic 的消费者
func (cm *ConsumerManager) GetConsumer(topic string) (*TopicConsumer, bool) {
	consumer, exists := cm.consumers[topic]
	return consumer, exists
}
