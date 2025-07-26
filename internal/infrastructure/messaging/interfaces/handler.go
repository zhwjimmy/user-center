package interfaces

import (
	"context"

	"github.com/IBM/sarama"
)

// Handler 单个 Topic 的处理器接口
type Handler interface {
	Handle(ctx context.Context, message *sarama.ConsumerMessage) error
	GetTopicName() string
	GetConsumerGroup() string
}
