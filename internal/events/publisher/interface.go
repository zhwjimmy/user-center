package publisher

import (
	"context"

	"github.com/zhwjimmy/user-center/internal/events/types"
)

// EventPublisher 事件发布接口
type EventPublisher interface {
	// 用户相关事件
	PublishUserRegistered(ctx context.Context, event *types.UserRegisteredEvent) error
	PublishUserLoggedIn(ctx context.Context, event *types.UserLoggedInEvent) error
	PublishUserPasswordChanged(ctx context.Context, event *types.UserPasswordChangedEvent) error
	PublishUserStatusChanged(ctx context.Context, event *types.UserStatusChangedEvent) error
	PublishUserDeleted(ctx context.Context, event *types.UserDeletedEvent) error
	PublishUserUpdated(ctx context.Context, event *types.UserUpdatedEvent) error

	// 关闭发布器
	Close() error
}
