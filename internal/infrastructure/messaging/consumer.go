package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/IBM/sarama"
	"github.com/zhwjimmy/user-center/internal/events/handlers"
	"github.com/zhwjimmy/user-center/internal/events/types"
	"go.uber.org/zap"
)

// MessageHandler 消息处理器接口
type MessageHandler interface {
	HandleMessage(ctx context.Context, message *sarama.ConsumerMessage) error
}

// EventHandlers 事件处理器集合
type EventHandlers struct {
	UserRegistered      *handlers.UserRegisteredHandler
	UserLoggedIn        *handlers.UserLoggedInHandler
	UserPasswordChanged *handlers.UserPasswordChangedHandler
	UserStatusChanged   *handlers.UserStatusChangedHandler
	UserDeleted         *handlers.UserDeletedHandler
	UserUpdated         *handlers.UserUpdatedHandler
}

// kafkaConsumer Kafka消费者实现
type kafkaConsumer struct {
	consumer sarama.ConsumerGroup
	config   *KafkaClientConfig
	handlers *EventHandlers
	logger   *zap.Logger
	topics   []string
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewKafkaConsumer 创建Kafka消费者
func NewKafkaConsumer(cfg *KafkaClientConfig, handlers *EventHandlers, logger *zap.Logger) (Consumer, error) {
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

	kc := &kafkaConsumer{
		consumer: consumer,
		config:   cfg,
		handlers: handlers,
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
func (c *kafkaConsumer) Start(ctx context.Context) error {
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
func (c *kafkaConsumer) Stop() error {
	c.logger.Info("Stopping Kafka consumer")
	if c.cancel != nil {
		c.cancel()
	}
	c.wg.Wait()
	return c.consumer.Close()
}

// Setup 实现 sarama.ConsumerGroupHandler 接口
func (c *kafkaConsumer) Setup(sarama.ConsumerGroupSession) error {
	c.logger.Info("Consumer group session setup")
	return nil
}

// Cleanup 实现 sarama.ConsumerGroupHandler 接口
func (c *kafkaConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	c.logger.Info("Consumer group session cleanup")
	return nil
}

// ConsumeClaim 实现 sarama.ConsumerGroupHandler 接口
func (c *kafkaConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
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
			if err := c.handleMessage(session.Context(), message); err != nil {
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

// handleMessage 处理消息
func (c *kafkaConsumer) handleMessage(ctx context.Context, message *sarama.ConsumerMessage) error {
	// 获取事件类型
	eventType := c.getEventType(message.Headers)

	c.logger.Debug("Processing message",
		zap.String("topic", message.Topic),
		zap.String("event_type", eventType),
		zap.Int32("partition", message.Partition),
		zap.Int64("offset", message.Offset),
	)

	switch types.EventType(eventType) {
	case types.UserRegistered:
		var userEvent types.UserRegisteredEvent
		if err := json.Unmarshal(message.Value, &userEvent); err != nil {
			return fmt.Errorf("failed to unmarshal user registered event: %w", err)
		}
		return c.handlers.UserRegistered.Handle(ctx, &userEvent)

	case types.UserLoggedIn:
		var userEvent types.UserLoggedInEvent
		if err := json.Unmarshal(message.Value, &userEvent); err != nil {
			return fmt.Errorf("failed to unmarshal user logged in event: %w", err)
		}
		return c.handlers.UserLoggedIn.Handle(ctx, &userEvent)

	case types.UserPasswordChanged:
		var userEvent types.UserPasswordChangedEvent
		if err := json.Unmarshal(message.Value, &userEvent); err != nil {
			return fmt.Errorf("failed to unmarshal user password changed event: %w", err)
		}
		return c.handlers.UserPasswordChanged.Handle(ctx, &userEvent)

	case types.UserStatusChanged:
		var userEvent types.UserStatusChangedEvent
		if err := json.Unmarshal(message.Value, &userEvent); err != nil {
			return fmt.Errorf("failed to unmarshal user status changed event: %w", err)
		}
		return c.handlers.UserStatusChanged.Handle(ctx, &userEvent)

	case types.UserDeleted:
		var userEvent types.UserDeletedEvent
		if err := json.Unmarshal(message.Value, &userEvent); err != nil {
			return fmt.Errorf("failed to unmarshal user deleted event: %w", err)
		}
		return c.handlers.UserDeleted.Handle(ctx, &userEvent)

	case types.UserUpdated:
		var userEvent types.UserUpdatedEvent
		if err := json.Unmarshal(message.Value, &userEvent); err != nil {
			return fmt.Errorf("failed to unmarshal user updated event: %w", err)
		}
		return c.handlers.UserUpdated.Handle(ctx, &userEvent)

	default:
		c.logger.Warn("Unknown event type", zap.String("event_type", eventType))
		return nil // 忽略未知事件类型
	}
}

// getEventType 从消息头获取事件类型
func (c *kafkaConsumer) getEventType(headers []*sarama.RecordHeader) string {
	for _, header := range headers {
		if string(header.Key) == "event_type" {
			return string(header.Value)
		}
	}
	return ""
}
