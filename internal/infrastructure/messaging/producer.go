package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/zhwjimmy/user-center/internal/kafka/event"
	"go.uber.org/zap"
)

// kafkaProducer Kafka生产者实现
type kafkaProducer struct {
	producer sarama.AsyncProducer
	config   *KafkaClientConfig
	logger   *zap.Logger
	wg       sync.WaitGroup
	closed   chan struct{}
}

// NewKafkaProducer 创建Kafka生产者
func NewKafkaProducer(cfg *KafkaClientConfig, logger *zap.Logger) (Producer, error) {
	producerConfig := cfg.NewProducerConfig()

	producer, err := sarama.NewAsyncProducer(cfg.Brokers, producerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer: %w", err)
	}

	kp := &kafkaProducer{
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
func (p *kafkaProducer) PublishUserEvent(ctx context.Context, event interface{}) error {
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
func (p *kafkaProducer) PublishUserEventAsync(ctx context.Context, event interface{}) error {
	message, err := p.createMessage(event)
	if err != nil {
		return err
	}

	// 异步发送
	select {
	case p.producer.Input() <- message:
		p.logger.Debug("Message sent to producer queue",
			zap.String("topic", message.Topic),
		)
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(5 * time.Second):
		return fmt.Errorf("timeout sending message to producer queue")
	}
}

// createMessage 创建Kafka消息
func (p *kafkaProducer) createMessage(eventData interface{}) (*sarama.ProducerMessage, error) {
	// 序列化事件数据
	jsonData, err := json.Marshal(eventData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal event data: %w", err)
	}

	// 确定主题名称
	var topic string
	switch eventData.(type) {
	case *event.UserRegisteredEvent:
		topic = p.config.GetTopicName("user_registered")
	case *event.UserLoggedInEvent:
		topic = p.config.GetTopicName("user_logged_in")
	case *event.UserPasswordChangedEvent:
		topic = p.config.GetTopicName("user_password_changed")
	case *event.UserStatusChangedEvent:
		topic = p.config.GetTopicName("user_status_changed")
	case *event.UserDeletedEvent:
		topic = p.config.GetTopicName("user_deleted")
	case *event.UserUpdatedEvent:
		topic = p.config.GetTopicName("user_updated")
	default:
		return nil, fmt.Errorf("unsupported event type: %T", eventData)
	}

	// 创建消息
	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(jsonData),
		Key:   sarama.StringEncoder(fmt.Sprintf("%d", time.Now().UnixNano())),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("content-type"),
				Value: []byte("application/json"),
			},
			{
				Key:   []byte("timestamp"),
				Value: []byte(time.Now().Format(time.RFC3339)),
			},
		},
	}

	return message, nil
}

// handleSuccesses 处理成功发送的消息
func (p *kafkaProducer) handleSuccesses() {
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

// handleErrors 处理发送失败的消息
func (p *kafkaProducer) handleErrors() {
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
func (p *kafkaProducer) Close() error {
	p.logger.Info("Closing Kafka producer")
	close(p.closed)
	p.producer.AsyncClose()
	p.wg.Wait()
	return nil
}
