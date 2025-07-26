// Package messaging 提供消息队列基础设施
//
// 架构层次：
// - interfaces/     : 接口定义层
// - kafka/         : Kafka 实现层
// - factory/       : 工厂层
package messaging

// 重新导出接口，保持向后兼容
import (
	"github.com/zhwjimmy/user-center/internal/infrastructure/messaging/factory"
	"github.com/zhwjimmy/user-center/internal/infrastructure/messaging/interfaces"
	"github.com/zhwjimmy/user-center/internal/infrastructure/messaging/kafka"
)

// 类型别名，保持向后兼容
type (
	Service        = interfaces.Service
	Producer       = interfaces.Producer
	Consumer       = interfaces.Consumer
	Handler        = interfaces.Handler
	Event          = interfaces.Event
	HandlerFactory = interfaces.HandlerFactory
	HandlerConfig  = interfaces.HandlerConfig
)

// 重新导出主要函数
var (
	NewKafkaService      = kafka.NewKafkaService
	NewKafkaProducer     = kafka.NewKafkaProducer
	NewKafkaConsumer     = kafka.NewKafkaConsumer
	NewKafkaClientConfig = kafka.NewKafkaClientConfig
)

// 重新导出工厂函数
func NewDefaultHandlerFactory() interfaces.HandlerFactory {
	return factory.NewDefaultHandlerFactory()
}
