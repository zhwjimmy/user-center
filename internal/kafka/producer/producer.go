package producer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/zhwjimmy/user-center/internal/kafka/config"
	"github.com/zhwjimmy/user-center/internal/kafka/event"
	"go.uber.org/zap"
)

// Producer Kafka生产者接口
type Producer interface {
	PublishUserEvent(ctx context.Context, event interface{}) error
	PublishUserEventAsync(ctx context.Context, event interface{}) error
	Close() error
}

// KafkaProducer Kafka生产者实现
type KafkaProducer struct {
	producer sarama.AsyncProducer
	config   *config.KafkaClientConfig
	logger   *zap.Logger
	wg       sync.WaitGroup
	closed   chan struct{}
}

// NewKafkaProducer 创建Kafka生产者
func NewKafkaProducer(cfg *config.KafkaClientConfig, logger *zap.Logger) (Producer, error) {
	producerConfig := cfg.NewProducerConfig()

	producer, err := sarama.NewAsyncProducer(cfg.Brokers, producerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer: %w", err)
	}

	kp := &KafkaProducer{
		producer: producer,
		config:   cfg,
		logger:   logger,
		closed:   make(chan struct{}),
	}

	// 启动错误和成功处理协程
	kp.wg.Add(2)
	go kp.handleSuccesses()
	go kp.handleErrors()

	logger.Info("Kafka producer created successfully",
		zap.Strings("brokers", cfg.Brokers),
		zap.String("group_id", cfg.GroupID),
	)

	return kp, nil
}

// PublishUserEvent 同步发布用户事件
func (p *KafkaProducer) PublishUserEvent(ctx context.Context, event interface{}) error {
	message, err := p.createMessage(event)
	if err != nil {
		return err
	}

	// 使用同步方式发送
	select {
	case p.producer.Input() <- message:
		// 等待确认
		select {
		case success := <-p.producer.Successes():
			p.logger.Debug("Message published successfully",
				zap.String("topic", success.Topic),
				zap.Int32("partition", success.Partition),
				zap.Int64("offset", success.Offset),
			)
			return nil
		case err := <-p.producer.Errors():
			p.logger.Error("Failed to publish message",
				zap.String("topic", err.Msg.Topic),
				zap.Error(err.Err),
			)
			return err.Err
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(30 * time.Second):
			return fmt.Errorf("timeout publishing message")
		}
	case <-ctx.Done():
		return ctx.Err()
	}
}

// PublishUserEventAsync 异步发布用户事件
func (p *KafkaProducer) PublishUserEventAsync(ctx context.Context, event interface{}) error {
	message, err := p.createMessage(event)
	if err != nil {
		return err
	}

	select {
	case p.producer.Input() <- message:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
		return fmt.Errorf("producer input channel is full")
	}
}

// createMessage 创建Kafka消息
func (p *KafkaProducer) createMessage(eventData interface{}) (*sarama.ProducerMessage, error) {
	var (
		topic   string
		key     string
		value   []byte
		headers []sarama.RecordHeader
	)

	switch e := eventData.(type) {
	case *event.UserRegisteredEvent:
		topic = p.config.GetTopicName("user_events")
		key = e.UserID
		var err error
		value, err = e.ToJSON()
		if err != nil {
			return nil, fmt.Errorf("failed to marshal user registered event: %w", err)
		}
		headers = []sarama.RecordHeader{
			{Key: []byte("event_type"), Value: []byte(e.Type)},
			{Key: []byte("request_id"), Value: []byte(e.RequestID)},
		}

	case *event.UserLoggedInEvent:
		topic = p.config.GetTopicName("user_events")
		key = e.UserID
		var err error
		value, err = e.ToJSON()
		if err != nil {
			return nil, fmt.Errorf("failed to marshal user logged in event: %w", err)
		}
		headers = []sarama.RecordHeader{
			{Key: []byte("event_type"), Value: []byte(e.Type)},
			{Key: []byte("request_id"), Value: []byte(e.RequestID)},
		}

	case *event.UserPasswordChangedEvent:
		topic = p.config.GetTopicName("user_events")
		key = e.UserID
		var err error
		value, err = e.ToJSON()
		if err != nil {
			return nil, fmt.Errorf("failed to marshal user password changed event: %w", err)
		}
		headers = []sarama.RecordHeader{
			{Key: []byte("event_type"), Value: []byte(e.Type)},
			{Key: []byte("request_id"), Value: []byte(e.RequestID)},
		}

	case *event.UserStatusChangedEvent:
		topic = p.config.GetTopicName("user_events")
		key = e.UserID
		var err error
		value, err = e.ToJSON()
		if err != nil {
			return nil, fmt.Errorf("failed to marshal user status changed event: %w", err)
		}
		headers = []sarama.RecordHeader{
			{Key: []byte("event_type"), Value: []byte(e.Type)},
			{Key: []byte("request_id"), Value: []byte(e.RequestID)},
		}

	case *event.UserDeletedEvent:
		topic = p.config.GetTopicName("user_events")
		key = e.UserID
		var err error
		value, err = e.ToJSON()
		if err != nil {
			return nil, fmt.Errorf("failed to marshal user deleted event: %w", err)
		}
		headers = []sarama.RecordHeader{
			{Key: []byte("event_type"), Value: []byte(e.Type)},
			{Key: []byte("request_id"), Value: []byte(e.RequestID)},
		}

	case *event.UserUpdatedEvent:
		topic = p.config.GetTopicName("user_events")
		key = e.UserID
		var err error
		value, err = e.ToJSON()
		if err != nil {
			return nil, fmt.Errorf("failed to marshal user updated event: %w", err)
		}
		headers = []sarama.RecordHeader{
			{Key: []byte("event_type"), Value: []byte(e.Type)},
			{Key: []byte("request_id"), Value: []byte(e.RequestID)},
		}

	default:
		return nil, fmt.Errorf("unsupported event type: %T", eventData)
	}

	return &sarama.ProducerMessage{
		Topic:     topic,
		Key:       sarama.StringEncoder(key),
		Value:     sarama.ByteEncoder(value),
		Headers:   headers,
		Timestamp: time.Now(),
	}, nil
}

// handleSuccesses 处理成功消息
func (p *KafkaProducer) handleSuccesses() {
	defer p.wg.Done()

	for {
		select {
		case success := <-p.producer.Successes():
			p.logger.Debug("Message published successfully",
				zap.String("topic", success.Topic),
				zap.Int32("partition", success.Partition),
				zap.Int64("offset", success.Offset),
			)
		case <-p.closed:
			return
		}
	}
}

// handleErrors 处理错误消息
func (p *KafkaProducer) handleErrors() {
	defer p.wg.Done()

	for {
		select {
		case err := <-p.producer.Errors():
			p.logger.Error("Failed to publish message",
				zap.String("topic", err.Msg.Topic),
				zap.Error(err.Err),
			)
		case <-p.closed:
			return
		}
	}
}

// Close 关闭生产者
func (p *KafkaProducer) Close() error {
	close(p.closed)

	if err := p.producer.Close(); err != nil {
		p.logger.Error("Failed to close kafka producer", zap.Error(err))
		return err
	}

	p.wg.Wait()
	p.logger.Info("Kafka producer closed successfully")
	return nil
}
