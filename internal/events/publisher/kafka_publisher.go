package publisher

import (
	"context"

	"github.com/zhwjimmy/user-center/internal/events/types"
	"github.com/zhwjimmy/user-center/internal/infrastructure/messaging"
	"go.uber.org/zap"
)

// KafkaEventPublisher Kafka 事件发布实现
type KafkaEventPublisher struct {
	producer messaging.Producer
	logger   *zap.Logger
}

// NewKafkaEventPublisher 创建 Kafka 事件发布器
func NewKafkaEventPublisher(producer messaging.Producer, logger *zap.Logger) EventPublisher {
	return &KafkaEventPublisher{
		producer: producer,
		logger:   logger,
	}
}

// PublishUserRegistered 发布用户注册事件
func (p *KafkaEventPublisher) PublishUserRegistered(ctx context.Context, event *types.UserRegisteredEvent) error {
	p.logger.Debug("Publishing user registered event", zap.String("user_id", event.UserID))
	return p.producer.PublishUserEvent(ctx, event)
}

// PublishUserLoggedIn 发布用户登录事件
func (p *KafkaEventPublisher) PublishUserLoggedIn(ctx context.Context, event *types.UserLoggedInEvent) error {
	p.logger.Debug("Publishing user logged in event", zap.String("user_id", event.UserID))
	return p.producer.PublishUserEvent(ctx, event)
}

// PublishUserPasswordChanged 发布用户密码变更事件
func (p *KafkaEventPublisher) PublishUserPasswordChanged(ctx context.Context, event *types.UserPasswordChangedEvent) error {
	p.logger.Debug("Publishing user password changed event", zap.String("user_id", event.UserID))
	return p.producer.PublishUserEvent(ctx, event)
}

// PublishUserStatusChanged 发布用户状态变更事件
func (p *KafkaEventPublisher) PublishUserStatusChanged(ctx context.Context, event *types.UserStatusChangedEvent) error {
	p.logger.Debug("Publishing user status changed event", zap.String("user_id", event.UserID))
	return p.producer.PublishUserEvent(ctx, event)
}

// PublishUserDeleted 发布用户删除事件
func (p *KafkaEventPublisher) PublishUserDeleted(ctx context.Context, event *types.UserDeletedEvent) error {
	p.logger.Debug("Publishing user deleted event", zap.String("user_id", event.UserID))
	return p.producer.PublishUserEvent(ctx, event)
}

// PublishUserUpdated 发布用户更新事件
func (p *KafkaEventPublisher) PublishUserUpdated(ctx context.Context, event *types.UserUpdatedEvent) error {
	p.logger.Debug("Publishing user updated event", zap.String("user_id", event.UserID))
	return p.producer.PublishUserEvent(ctx, event)
}

// Close 关闭发布器
func (p *KafkaEventPublisher) Close() error {
	return p.producer.Close()
}
