package kafka

import (
	"context"
	"fmt"

	"github.com/zhwjimmy/user-center/internal/kafka/config"
	"github.com/zhwjimmy/user-center/internal/kafka/consumer"
	"github.com/zhwjimmy/user-center/internal/kafka/producer"
	"go.uber.org/zap"
)

// Service Kafka服务接口
type Service interface {
	GetProducer() producer.Producer
	GetConsumer() consumer.Consumer
	Start(ctx context.Context) error
	Stop() error
}

// KafkaService Kafka服务实现
type KafkaService struct {
	producer producer.Producer
	consumer consumer.Consumer
	logger   *zap.Logger
}

// NewKafkaService 创建Kafka服务
func NewKafkaService(cfg *config.KafkaClientConfig, logger *zap.Logger) (Service, error) {
	// 创建生产者
	prod, err := producer.NewKafkaProducer(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer: %w", err)
	}

	// 创建消息处理器
	handler := consumer.NewUserEventHandler(logger)

	// 创建消费者
	cons, err := consumer.NewKafkaConsumer(cfg, handler, logger)
	if err != nil {
		prod.Close() // 清理已创建的生产者
		return nil, fmt.Errorf("failed to create kafka consumer: %w", err)
	}

	return &KafkaService{
		producer: prod,
		consumer: cons,
		logger:   logger,
	}, nil
}

// GetProducer 获取生产者
func (s *KafkaService) GetProducer() producer.Producer {
	return s.producer
}

// GetConsumer 获取消费者
func (s *KafkaService) GetConsumer() consumer.Consumer {
	return s.consumer
}

// Start 启动Kafka服务
func (s *KafkaService) Start(ctx context.Context) error {
	s.logger.Info("Starting Kafka service")

	// 启动消费者
	if err := s.consumer.Start(ctx); err != nil {
		return fmt.Errorf("failed to start kafka consumer: %w", err)
	}

	s.logger.Info("Kafka service started successfully")
	return nil
}

// Stop 停止Kafka服务
func (s *KafkaService) Stop() error {
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
