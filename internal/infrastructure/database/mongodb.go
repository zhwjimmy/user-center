package database

import (
	"context"
	"fmt"
	"time"

	"github.com/zhwjimmy/user-center/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// mongoDB 实现 MongoDB 接口
type mongoDB struct {
	client   *mongo.Client
	database *mongo.Database
}

// 确保 mongoDB 实现了 MongoDB 接口
var _ MongoDB = (*mongoDB)(nil)

// NewMongoDB 创建新的 MongoDB 连接
func NewMongoDB(cfg *config.Config, logger *zap.Logger) (MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set client options
	clientOptions := options.Client().ApplyURI(cfg.Database.MongoDB.URI)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Test connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	database := client.Database(cfg.Database.MongoDB.Database)

	logger.Info("MongoDB connected successfully",
		zap.String("uri", cfg.Database.MongoDB.URI),
		zap.String("database", cfg.Database.MongoDB.Database),
	)

	return &mongoDB{
		client:   client,
		database: database,
	}, nil
}

// Close 关闭 MongoDB 连接
func (m *mongoDB) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}

// Health 检查 MongoDB 健康状态
func (m *mongoDB) Health(ctx context.Context) error {
	return m.client.Ping(ctx, nil)
}

// Collection 返回集合实例
func (m *mongoDB) Collection(name string) interface{} {
	return m.database.Collection(name)
}

// LogEntry 表示 MongoDB 中的日志条目
type LogEntry struct {
	ID        string                 `bson:"_id,omitempty"`
	Level     string                 `bson:"level"`
	Message   string                 `bson:"message"`
	Timestamp time.Time              `bson:"timestamp"`
	RequestID string                 `bson:"request_id,omitempty"`
	UserID    uint                   `bson:"user_id,omitempty"`
	Fields    map[string]interface{} `bson:"fields,omitempty"`
}

// UserSession 表示 MongoDB 中的用户会话
type UserSession struct {
	ID        string    `bson:"_id,omitempty"`
	UserID    uint      `bson:"user_id"`
	Token     string    `bson:"token"`
	IP        string    `bson:"ip"`
	UserAgent string    `bson:"user_agent"`
	CreatedAt time.Time `bson:"created_at"`
	ExpiresAt time.Time `bson:"expires_at"`
	IsActive  bool      `bson:"is_active"`
}

// AuditLog 表示 MongoDB 中的审计日志条目
type AuditLog struct {
	ID        string                 `bson:"_id,omitempty"`
	UserID    uint                   `bson:"user_id,omitempty"`
	Action    string                 `bson:"action"`
	Resource  string                 `bson:"resource"`
	Details   map[string]interface{} `bson:"details,omitempty"`
	IP        string                 `bson:"ip,omitempty"`
	UserAgent string                 `bson:"user_agent,omitempty"`
	Timestamp time.Time              `bson:"timestamp"`
	RequestID string                 `bson:"request_id,omitempty"`
}
