package messaging

import (
	"time"

	"github.com/IBM/sarama"
	"github.com/zhwjimmy/user-center/internal/config"
)

// KafkaClientConfig Kafka客户端配置
type KafkaClientConfig struct {
	Brokers       []string
	Topics        map[string]string
	GroupID       string
	RetryMax      int
	RetryBackoff  time.Duration
	BatchSize     int
	BatchTimeout  time.Duration
	FlushMessages int
	FlushBytes    int
	Compression   sarama.CompressionCodec
}

// NewKafkaClientConfig 创建Kafka客户端配置
func NewKafkaClientConfig(cfg *config.Config) *KafkaClientConfig {
	return &KafkaClientConfig{
		Brokers:       cfg.Kafka.Brokers,
		Topics:        cfg.Kafka.Topics,
		GroupID:       cfg.Kafka.GroupID,
		RetryMax:      3,
		RetryBackoff:  100 * time.Millisecond,
		BatchSize:     100,
		BatchTimeout:  10 * time.Millisecond,
		FlushMessages: 100,
		FlushBytes:    1024 * 1024, // 1MB
		Compression:   sarama.CompressionSnappy,
	}
}

// NewProducerConfig 创建生产者配置
func (c *KafkaClientConfig) NewProducerConfig() *sarama.Config {
	config := sarama.NewConfig()

	// 生产者配置
	config.Producer.RequiredAcks = sarama.WaitForAll // 等待所有副本确认
	config.Producer.Retry.Max = c.RetryMax
	config.Producer.Retry.Backoff = c.RetryBackoff
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Compression = c.Compression

	// 批处理配置
	config.Producer.Flush.Messages = c.FlushMessages
	config.Producer.Flush.Bytes = c.FlushBytes
	config.Producer.Flush.Frequency = c.BatchTimeout

	// 幂等性配置
	config.Producer.Idempotent = true
	config.Net.MaxOpenRequests = 1

	// 版本配置
	config.Version = sarama.V2_6_0_0

	return config
}

// NewConsumerConfig 创建消费者配置
func (c *KafkaClientConfig) NewConsumerConfig() *sarama.Config {
	config := sarama.NewConfig()

	// 消费者配置
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Group.Session.Timeout = 10 * time.Second
	config.Consumer.Group.Heartbeat.Interval = 3 * time.Second

	// 自动提交offset
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second

	// 版本配置
	config.Version = sarama.V2_6_0_0

	return config
}

// GetTopicName 获取主题名称
func (c *KafkaClientConfig) GetTopicName(key string) string {
	if topic, exists := c.Topics[key]; exists {
		return topic
	}
	return key
}
