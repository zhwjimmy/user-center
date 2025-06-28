package database

import (
	"context"
	"fmt"
	"time"

	"github.com/your-org/user-center/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// MongoDB represents MongoDB database connection
type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

// NewMongoDB creates a new MongoDB connection
func NewMongoDB(cfg *config.Config, logger *zap.Logger) (*MongoDB, error) {
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

	return &MongoDB{
		Client:   client,
		Database: database,
	}, nil
}

// Close closes the MongoDB connection
func (m *MongoDB) Close(ctx context.Context) error {
	return m.Client.Disconnect(ctx)
}

// Health checks the MongoDB health
func (m *MongoDB) Health(ctx context.Context) error {
	return m.Client.Ping(ctx, nil)
}

// Collection returns a collection instance
func (m *MongoDB) Collection(name string) *mongo.Collection {
	return m.Database.Collection(name)
}

// LogEntry represents a log entry in MongoDB
type LogEntry struct {
	ID        string                 `bson:"_id,omitempty"`
	Level     string                 `bson:"level"`
	Message   string                 `bson:"message"`
	Timestamp time.Time              `bson:"timestamp"`
	RequestID string                 `bson:"request_id,omitempty"`
	UserID    uint                   `bson:"user_id,omitempty"`
	Fields    map[string]interface{} `bson:"fields,omitempty"`
}

// UserSession represents a user session in MongoDB
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

// AuditLog represents an audit log entry in MongoDB
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
