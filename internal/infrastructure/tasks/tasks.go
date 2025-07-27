// Package tasks 提供任务队列基础设施
//
// 架构层次：
// - interfaces/     : 接口定义层
// - asynq/         : Asynq 实现层
// - factory/       : 工厂层
package tasks

// 重新导出接口，保持向后兼容
import (
	"github.com/zhwjimmy/user-center/internal/infrastructure/tasks/asynq"
	"github.com/zhwjimmy/user-center/internal/infrastructure/tasks/factory"
	"github.com/zhwjimmy/user-center/internal/infrastructure/tasks/interfaces"
)

// 类型别名，保持向后兼容
type (
	Service        = interfaces.Service
	Client         = interfaces.Client
	Server         = interfaces.Server
	Handler        = interfaces.Handler
	HandlerFactory = interfaces.HandlerFactory
)

// 重新导出主要函数
var (
	NewAsynqService      = asynq.NewAsynqService
	NewAsynqClient       = asynq.NewAsynqClient
	NewAsynqServer       = asynq.NewAsynqServer
	NewAsynqClientConfig = asynq.NewAsynqConfig
)

// 重新导出工厂函数
func NewDefaultHandlerFactory() interfaces.HandlerFactory {
	return factory.NewDefaultHandlerFactory()
}
