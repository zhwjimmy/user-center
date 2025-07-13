package producer

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"github.com/zhwjimmy/user-center/internal/kafka/config"
	"github.com/zhwjimmy/user-center/internal/kafka/event"
	"go.uber.org/zap/zaptest"
)

func TestKafkaProducer(t *testing.T) {
	// Skip integration test in CI environment or when Kafka is not available
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Check if Kafka is available
	if !isKafkaAvailable("localhost:9092") {
		t.Skip("skipping test: Kafka not available")
	}

	logger := zaptest.NewLogger(t)
	cfg := &config.KafkaClientConfig{
		Brokers: []string{"localhost:9092"},
		Topics: map[string]string{
			"user_events": "test.user.events",
		},
		GroupID:       "test-group",
		RetryMax:      3,
		RetryBackoff:  100 * time.Millisecond,
		BatchSize:     100,
		BatchTimeout:  10 * time.Millisecond,
		FlushMessages: 100,
		FlushBytes:    1024 * 1024, // 1MB
		Compression:   sarama.CompressionSnappy,
	}

	producer, err := NewKafkaProducer(cfg, logger)
	if err != nil {
		t.Skipf("skipping test: failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	ctx := context.Background()

	// 测试发布用户注册事件
	userEvent := &event.UserRegisteredEvent{
		BaseEvent: event.NewBaseEvent(
			event.UserRegistered,
			"test-source",
			"test-request-id",
			"test-user-id",
		),
		Username: "testuser",
		Email:    "test@example.com",
	}

	err = producer.PublishUserEventAsync(ctx, userEvent)
	assert.NoError(t, err)

	// 等待消息发送完成
	time.Sleep(100 * time.Millisecond)
}

// isKafkaAvailable checks if Kafka is available at the given address
func isKafkaAvailable(addr string) bool {
	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
