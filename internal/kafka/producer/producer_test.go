package producer

import (
	"context"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zhwjimmy/user-center/internal/kafka/config"
	"github.com/zhwjimmy/user-center/internal/kafka/event"
	"go.uber.org/zap/zaptest"
)

func TestKafkaProducer(t *testing.T) {
	// 跳过集成测试
	if testing.Short() {
		t.Skip("skipping integration test")
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
	require.NoError(t, err)
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
