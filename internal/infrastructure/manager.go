package infrastructure

import (
	"context"
	"fmt"
	"sync"

	"github.com/zhwjimmy/user-center/internal/config"
	"github.com/zhwjimmy/user-center/internal/infrastructure/cache"
	"github.com/zhwjimmy/user-center/internal/infrastructure/database"
	"github.com/zhwjimmy/user-center/internal/infrastructure/messaging"
	"go.uber.org/zap"
)

// Manager 基础设施管理器
type Manager struct {
	config *config.Config
	logger *zap.Logger

	// 数据库连接
	postgres database.PostgreSQL
	mongodb  database.MongoDB

	// 缓存连接
	redis cache.Cache

	// 消息队列
	kafka messaging.Service

	// 生命周期管理
	mu     sync.RWMutex
	closed bool
}

// NewManager 创建基础设施管理器
func NewManager(cfg *config.Config, logger *zap.Logger) (*Manager, error) {
	manager := &Manager{
		config: cfg,
		logger: logger,
	}

	// 初始化数据库连接
	if err := manager.initDatabases(); err != nil {
		return nil, fmt.Errorf("failed to initialize databases: %w", err)
	}

	// 初始化缓存连接
	if err := manager.initCache(); err != nil {
		return nil, fmt.Errorf("failed to initialize cache: %w", err)
	}

	// 初始化消息队列
	if err := manager.initMessaging(); err != nil {
		return nil, fmt.Errorf("failed to initialize messaging: %w", err)
	}

	logger.Info("Infrastructure manager initialized successfully")
	return manager, nil
}

// initDatabases 初始化数据库连接
func (m *Manager) initDatabases() error {
	// 初始化 PostgreSQL
	postgres, err := database.NewPostgreSQL(m.config, m.logger)
	if err != nil {
		return fmt.Errorf("failed to initialize PostgreSQL: %w", err)
	}
	m.postgres = postgres

	// 初始化 MongoDB
	mongodb, err := database.NewMongoDB(m.config, m.logger)
	if err != nil {
		return fmt.Errorf("failed to initialize MongoDB: %w", err)
	}
	m.mongodb = mongodb

	return nil
}

// initCache 初始化缓存连接
func (m *Manager) initCache() error {
	redis, err := cache.NewRedis(m.config, m.logger)
	if err != nil {
		return fmt.Errorf("failed to initialize Redis: %w", err)
	}
	m.redis = redis

	return nil
}

// initMessaging 初始化消息队列
func (m *Manager) initMessaging() error {
	// 创建 Kafka 配置
	kafkaConfig := messaging.NewKafkaClientConfig(m.config)

	kafka, err := messaging.NewKafkaService(kafkaConfig, m.logger)
	if err != nil {
		return fmt.Errorf("failed to initialize Kafka: %w", err)
	}
	m.kafka = kafka

	return nil
}

// GetPostgreSQL 获取 PostgreSQL 连接
func (m *Manager) GetPostgreSQL() database.PostgreSQL {
	return m.postgres
}

// GetMongoDB 获取 MongoDB 连接
func (m *Manager) GetMongoDB() database.MongoDB {
	return m.mongodb
}

// GetRedis 获取 Redis 连接
func (m *Manager) GetRedis() cache.Cache {
	return m.redis
}

// GetKafka 获取 Kafka 服务
func (m *Manager) GetKafka() messaging.Service {
	return m.kafka
}

// Start 启动所有服务
func (m *Manager) Start(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return fmt.Errorf("manager is already closed")
	}

	// 启动 Kafka 服务
	if err := m.kafka.Start(ctx); err != nil {
		return fmt.Errorf("failed to start kafka service: %w", err)
	}

	m.logger.Info("All infrastructure services started successfully")
	return nil
}

// Stop 停止所有服务
func (m *Manager) Stop(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return nil
	}

	var errors []error

	// 停止 Kafka 服务
	if err := m.kafka.Stop(); err != nil {
		errors = append(errors, fmt.Errorf("failed to stop kafka: %w", err))
	}

	// 关闭 Redis 连接
	if err := m.redis.Close(); err != nil {
		errors = append(errors, fmt.Errorf("failed to close redis: %w", err))
	}

	// 关闭 MongoDB 连接
	if err := m.mongodb.Close(ctx); err != nil {
		errors = append(errors, fmt.Errorf("failed to close mongodb: %w", err))
	}

	// 关闭 PostgreSQL 连接
	if err := m.postgres.Close(); err != nil {
		errors = append(errors, fmt.Errorf("failed to close postgresql: %w", err))
	}

	m.closed = true
	m.logger.Info("Infrastructure manager stopped")

	if len(errors) > 0 {
		return fmt.Errorf("errors during shutdown: %v", errors)
	}

	return nil
}

// IsClosed 检查管理器是否已关闭
func (m *Manager) IsClosed() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.closed
}
