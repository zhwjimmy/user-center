package consumer

import (
	"context"
	"fmt"
	"sync"

	"github.com/IBM/sarama"
	"github.com/zhwjimmy/user-center/internal/kafka/config"
	"github.com/zhwjimmy/user-center/internal/kafka/event"
	"go.uber.org/zap"
)

// MessageHandler 消息处理器接口
type MessageHandler interface {
	HandleUserRegistered(ctx context.Context, event *event.UserRegisteredEvent) error
	HandleUserLoggedIn(ctx context.Context, event *event.UserLoggedInEvent) error
	HandleUserPasswordChanged(ctx context.Context, event *event.UserPasswordChangedEvent) error
	HandleUserStatusChanged(ctx context.Context, event *event.UserStatusChangedEvent) error
	HandleUserDeleted(ctx context.Context, event *event.UserDeletedEvent) error
	HandleUserUpdated(ctx context.Context, event *event.UserUpdatedEvent) error
}

// Consumer Kafka消费者接口
type Consumer interface {
	Start(ctx context.Context) error
	Stop() error
}

// KafkaConsumer Kafka消费者实现
type KafkaConsumer struct {
	consumerGroup sarama.ConsumerGroup
	config        *config.KafkaClientConfig
	handler       MessageHandler
	logger        *zap.Logger
	wg            sync.WaitGroup
	cancel        context.CancelFunc
}

// NewKafkaConsumer 创建Kafka消费者
func NewKafkaConsumer(cfg *config.KafkaClientConfig, handler MessageHandler, logger *zap.Logger) (Consumer, error) {
	consumerConfig := cfg.NewConsumerConfig()

	consumerGroup, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.GroupID, consumerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka consumer group: %w", err)
	}

	consumer := &KafkaConsumer{
		consumerGroup: consumerGroup,
		config:        cfg,
		handler:       handler,
		logger:        logger,
	}

	logger.Info("Kafka consumer created successfully",
		zap.Strings("brokers", cfg.Brokers),
		zap.String("group_id", cfg.GroupID),
	)

	return consumer, nil
}

// Start 启动消费者
func (c *KafkaConsumer) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	c.cancel = cancel

	topics := []string{c.config.GetTopicName("user_events")}

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()

		for {
			select {
			case <-ctx.Done():
				c.logger.Info("Consumer context cancelled")
				return
			default:
				if err := c.consumerGroup.Consume(ctx, topics, c); err != nil {
					c.logger.Error("Error consuming messages", zap.Error(err))
					return
				}
			}
		}
	}()

	// 处理错误
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()

		for {
			select {
			case err := <-c.consumerGroup.Errors():
				c.logger.Error("Consumer group error", zap.Error(err))
			case <-ctx.Done():
				return
			}
		}
	}()

	c.logger.Info("Kafka consumer started", zap.Strings("topics", topics))
	return nil
}

// Stop 停止消费者
func (c *KafkaConsumer) Stop() error {
	if c.cancel != nil {
		c.cancel()
	}

	if err := c.consumerGroup.Close(); err != nil {
		c.logger.Error("Failed to close consumer group", zap.Error(err))
		return err
	}

	c.wg.Wait()
	c.logger.Info("Kafka consumer stopped successfully")
	return nil
}

// Setup 消费者组设置
func (c *KafkaConsumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup 消费者组清理
func (c *KafkaConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim 消费消息
func (c *KafkaConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				return nil
			}

			if err := c.processMessage(session.Context(), message); err != nil {
				c.logger.Error("Failed to process message",
					zap.String("topic", message.Topic),
					zap.Int32("partition", message.Partition),
					zap.Int64("offset", message.Offset),
					zap.Error(err),
				)
				// 继续处理下一条消息，不中断消费
				continue
			}

			// 标记消息已处理
			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}

// processMessage 处理消息
func (c *KafkaConsumer) processMessage(ctx context.Context, message *sarama.ConsumerMessage) error {
	// 获取事件类型
	eventType := c.getEventType(message.Headers)

	c.logger.Debug("Processing message",
		zap.String("topic", message.Topic),
		zap.String("event_type", eventType),
		zap.Int32("partition", message.Partition),
		zap.Int64("offset", message.Offset),
	)

	switch event.EventType(eventType) {
	case event.UserRegistered:
		var userEvent event.UserRegisteredEvent
		if err := userEvent.FromJSON(message.Value); err != nil {
			return fmt.Errorf("failed to unmarshal user registered event: %w", err)
		}
		return c.handler.HandleUserRegistered(ctx, &userEvent)

	case event.UserLoggedIn:
		var userEvent event.UserLoggedInEvent
		if err := userEvent.FromJSON(message.Value); err != nil {
			return fmt.Errorf("failed to unmarshal user logged in event: %w", err)
		}
		return c.handler.HandleUserLoggedIn(ctx, &userEvent)

	case event.UserPasswordChanged:
		var userEvent event.UserPasswordChangedEvent
		if err := userEvent.FromJSON(message.Value); err != nil {
			return fmt.Errorf("failed to unmarshal user password changed event: %w", err)
		}
		return c.handler.HandleUserPasswordChanged(ctx, &userEvent)

	case event.UserStatusChanged:
		var userEvent event.UserStatusChangedEvent
		if err := userEvent.FromJSON(message.Value); err != nil {
			return fmt.Errorf("failed to unmarshal user status changed event: %w", err)
		}
		return c.handler.HandleUserStatusChanged(ctx, &userEvent)

	case event.UserDeleted:
		var userEvent event.UserDeletedEvent
		if err := userEvent.FromJSON(message.Value); err != nil {
			return fmt.Errorf("failed to unmarshal user deleted event: %w", err)
		}
		return c.handler.HandleUserDeleted(ctx, &userEvent)

	case event.UserUpdated:
		var userEvent event.UserUpdatedEvent
		if err := userEvent.FromJSON(message.Value); err != nil {
			return fmt.Errorf("failed to unmarshal user updated event: %w", err)
		}
		return c.handler.HandleUserUpdated(ctx, &userEvent)

	default:
		c.logger.Warn("Unknown event type", zap.String("event_type", eventType))
		return nil // 忽略未知事件类型
	}
}

// getEventType 从消息头获取事件类型
func (c *KafkaConsumer) getEventType(headers []*sarama.RecordHeader) string {
	for _, header := range headers {
		if string(header.Key) == "event_type" {
			return string(header.Value)
		}
	}
	return ""
}
