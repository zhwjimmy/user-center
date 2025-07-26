package consumer

import (
	"context"
	"fmt"
	"sync"

	"github.com/IBM/sarama"
	"github.com/zhwjimmy/user-center/internal/infrastructure/messaging"
	"go.uber.org/zap"
)

// Consumer Kafka消费者接口
type Consumer interface {
	Start(ctx context.Context) error
	Stop() error
}

// KafkaConsumer Kafka消费者实现
type KafkaConsumer struct {
	consumer sarama.ConsumerGroup
	config   *messaging.KafkaClientConfig
	handler  MessageHandler
	logger   *zap.Logger
	topics   []string
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
}

// MessageHandler 消息处理器接口
type MessageHandler interface {
	HandleMessage(ctx context.Context, message *sarama.ConsumerMessage) error
}

// NewKafkaConsumer 创建Kafka消费者
func NewKafkaConsumer(cfg *messaging.KafkaClientConfig, handler MessageHandler, logger *zap.Logger) (Consumer, error) {
	consumerConfig := cfg.NewConsumerConfig()

	consumer, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.GroupID, consumerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka consumer: %w", err)
	}

	// 确定要消费的主题
	topics := []string{
		cfg.GetTopicName("user_registered"),
		cfg.GetTopicName("user_logged_in"),
		cfg.GetTopicName("user_password_changed"),
		cfg.GetTopicName("user_status_changed"),
		cfg.GetTopicName("user_deleted"),
		cfg.GetTopicName("user_updated"),
	}

	kc := &KafkaConsumer{
		consumer: consumer,
		config:   cfg,
		handler:  handler,
		logger:   logger,
		topics:   topics,
	}

	logger.Info("Kafka consumer created successfully",
		zap.Strings("brokers", cfg.Brokers),
		zap.String("group_id", cfg.GroupID),
		zap.Strings("topics", topics),
	)

	return kc, nil
}

// Start 启动消费者
func (c *KafkaConsumer) Start(ctx context.Context) error {
	c.ctx, c.cancel = context.WithCancel(ctx)

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for {
			select {
			case <-c.ctx.Done():
				return
			default:
				if err := c.consumer.Consume(c.ctx, c.topics, c); err != nil {
					c.logger.Error("Error from consumer", zap.Error(err))
				}
			}
		}
	}()

	c.logger.Info("Kafka consumer started")
	return nil
}

// Stop 停止消费者
func (c *KafkaConsumer) Stop() error {
	c.logger.Info("Stopping Kafka consumer")
	if c.cancel != nil {
		c.cancel()
	}
	c.wg.Wait()
	return c.consumer.Close()
}

// Setup 实现 sarama.ConsumerGroupHandler 接口
func (c *KafkaConsumer) Setup(sarama.ConsumerGroupSession) error {
	c.logger.Info("Consumer group session setup")
	return nil
}

// Cleanup 实现 sarama.ConsumerGroupHandler 接口
func (c *KafkaConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	c.logger.Info("Consumer group session cleanup")
	return nil
}

// ConsumeClaim 实现 sarama.ConsumerGroupHandler 接口
func (c *KafkaConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				continue
			}

			c.logger.Debug("Received message",
				zap.String("topic", message.Topic),
				zap.Int32("partition", message.Partition),
				zap.Int64("offset", message.Offset),
			)

			// 处理消息
			if err := c.handler.HandleMessage(session.Context(), message); err != nil {
				c.logger.Error("Failed to handle message",
					zap.String("topic", message.Topic),
					zap.Int32("partition", message.Partition),
					zap.Int64("offset", message.Offset),
					zap.Error(err),
				)
				// 不标记消息为已处理，让 Kafka 重新投递
				continue
			}

			// 标记消息为已处理
			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}
