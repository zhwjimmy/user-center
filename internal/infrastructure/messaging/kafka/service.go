package kafka

import (
	"context"
	"fmt"

	"github.com/zhwjimmy/user-center/internal/infrastructure/messaging/interfaces"
	"go.uber.org/zap"
)

// kafkaService 实现 Service 接口
type kafkaService struct {
	producer interfaces.Producer
	consumer interfaces.Consumer
	logger   *zap.Logger
}

// 确保 kafkaService 实现了 Service 接口
var _ interfaces.Service = (*kafkaService)(nil)

// NewKafkaService 创建 Kafka 服务
func NewKafkaService(cfg *KafkaClientConfig, handlerFactory interfaces.HandlerFactory, logger *zap.Logger) (interfaces.Service, error) {
	// 创建生产者
	prod, err := NewKafkaProducer(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer: %w", err)
	}

	// 通过工厂创建 handlers
	handlers, err := handlerFactory.CreateHandlers(cfg, logger)
	if err != nil {
		prod.Close() // 清理已创建的生产者
		return nil, fmt.Errorf("failed to create handlers: %w", err)
	}

	// 创建消费者
	cons, err := NewKafkaConsumer(cfg, handlers, logger)
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
func (s *kafkaService) GetProducer() interfaces.Producer {
	return s.producer
}

// GetConsumer 获取消费者
func (s *kafkaService) GetConsumer() interfaces.Consumer {
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
