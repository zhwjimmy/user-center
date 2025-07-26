package messaging

import (
	"context"
	"fmt"

	"github.com/zhwjimmy/user-center/internal/events/handlers"
	"go.uber.org/zap"
)

// kafkaService 实现 Service 接口
type kafkaService struct {
	producer Producer
	consumer Consumer
	logger   *zap.Logger
}

// 确保 kafkaService 实现了 Service 接口
var _ Service = (*kafkaService)(nil)

// NewKafkaService 创建 Kafka 服务
func NewKafkaService(cfg *KafkaClientConfig, logger *zap.Logger) (Service, error) {
	// 创建生产者
	prod, err := NewKafkaProducer(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer: %w", err)
	}

	// 创建事件处理器
	eventHandlers := &EventHandlers{
		UserRegistered:      handlers.NewUserRegisteredHandler(logger),
		UserLoggedIn:        handlers.NewUserLoggedInHandler(logger),
		UserPasswordChanged: handlers.NewUserPasswordChangedHandler(logger),
		UserStatusChanged:   handlers.NewUserStatusChangedHandler(logger),
		UserDeleted:         handlers.NewUserDeletedHandler(logger),
		UserUpdated:         handlers.NewUserUpdatedHandler(logger),
	}

	// 创建消费者
	cons, err := NewKafkaConsumer(cfg, eventHandlers, logger)
	if err != nil {
		prod.Close() // 清理已创建的生产者
		return nil, fmt.Errorf("failed to create kafka consumer: %w", err)
	}

	return &kafkaService{
		producer: prod,
		consumer: cons,
		logger:   logger,
	}, nil
}

// GetProducer 获取生产者
func (s *kafkaService) GetProducer() Producer {
	return s.producer
}

// GetConsumer 获取消费者
func (s *kafkaService) GetConsumer() Consumer {
	return s.consumer
}

// Start 启动 Kafka 服务
func (s *kafkaService) Start(ctx context.Context) error {
	s.logger.Info("Starting Kafka service")

	// 启动消费者
	if err := s.consumer.Start(ctx); err != nil {
		return fmt.Errorf("failed to start kafka consumer: %w", err)
	}

	s.logger.Info("Kafka service started successfully")
	return nil
}

// Stop 停止 Kafka 服务
func (s *kafkaService) Stop() error {
	s.logger.Info("Stopping Kafka service")

	// 停止消费者
	if err := s.consumer.Stop(); err != nil {
		s.logger.Error("Failed to stop kafka consumer", zap.Error(err))
	}

	// 停止生产者
	if err := s.producer.Close(); err != nil {
		s.logger.Error("Failed to close kafka producer", zap.Error(err))
	}

	s.logger.Info("Kafka service stopped successfully")
	return nil
}
