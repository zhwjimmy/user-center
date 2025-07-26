package kafka

import (
	"context"
	"fmt"
	"sync"

	"github.com/IBM/sarama"
	"github.com/zhwjimmy/user-center/internal/infrastructure/messaging/interfaces"
	"go.uber.org/zap"
)

// TopicConsumer
type TopicConsumer struct {
	handler  interfaces.Handler
	consumer sarama.ConsumerGroup
	config   *KafkaClientConfig
	logger   *zap.Logger
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewTopicConsumer 创建单个 Topic 的消费者
func NewTopicConsumer(cfg *KafkaClientConfig, handler interfaces.Handler, logger *zap.Logger) (*TopicConsumer, error) {
	consumerConfig := cfg.newConsumerConfig()

	consumer, err := sarama.NewConsumerGroup(cfg.Brokers, handler.GetConsumerGroup(), consumerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer for topic %s: %w", handler.GetTopicName(), err)
	}

	return &TopicConsumer{
		handler:  handler,
		consumer: consumer,
		config:   cfg,
		logger:   logger,
	}, nil
}

// Start 启动单个 Topic 消费者
func (tc *TopicConsumer) Start(ctx context.Context) error {
	tc.ctx, tc.cancel = context.WithCancel(ctx)

	tc.wg.Add(1)
	go func() {
		defer tc.wg.Done()
		for {
			select {
			case <-tc.ctx.Done():
				return
			default:
				if err := tc.consumer.Consume(tc.ctx, []string{tc.handler.GetTopicName()}, tc); err != nil {
					tc.logger.Error("Error from consumer",
						zap.String("topic", tc.handler.GetTopicName()),
						zap.String("group", tc.handler.GetConsumerGroup()),
						zap.Error(err))
				}
			}
		}
	}()

	tc.logger.Info("Topic consumer started",
		zap.String("topic", tc.handler.GetTopicName()),
		zap.String("group", tc.handler.GetConsumerGroup()))
	return nil
}

// Stop 停止单个 Topic 消费者
func (tc *TopicConsumer) Stop() error {
	tc.logger.Info("Stopping topic consumer",
		zap.String("topic", tc.handler.GetTopicName()),
		zap.String("group", tc.handler.GetConsumerGroup()))

	if tc.cancel != nil {
		tc.cancel()
	}
	tc.wg.Wait()
	return tc.consumer.Close()
}

// Setup 实现 sarama.ConsumerGroupHandler 接口
func (tc *TopicConsumer) Setup(sarama.ConsumerGroupSession) error {
	tc.logger.Info("Topic consumer group session setup",
		zap.String("topic", tc.handler.GetTopicName()),
		zap.String("group", tc.handler.GetConsumerGroup()))
	return nil
}

// Cleanup 实现 sarama.ConsumerGroupHandler 接口
func (tc *TopicConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	tc.logger.Info("Topic consumer group session cleanup",
		zap.String("topic", tc.handler.GetTopicName()),
		zap.String("group", tc.handler.GetConsumerGroup()))
	return nil
}

// ConsumeClaim 实现 sarama.ConsumerGroupHandler 接口
func (tc *TopicConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				continue
			}

			tc.logger.Debug("Received message",
				zap.String("topic", message.Topic),
				zap.Int32("partition", message.Partition),
				zap.Int64("offset", message.Offset),
			)

			// 直接调用 Handler 处理消息
			if err := tc.handler.Handle(session.Context(), message); err != nil {
				tc.logger.Error("Failed to handle message",
					zap.String("topic", message.Topic),
					zap.Int32("partition", message.Partition),
					zap.Int64("offset", message.Offset),
					zap.Error(err),
				)
				continue
			}

			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}
