package service

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/zhwjimmy/user-center/internal/events/publisher"
	"github.com/zhwjimmy/user-center/internal/events/types"
	"go.uber.org/zap"
)

// EventService 事件服务
type EventService struct {
	publisher publisher.EventPublisher
	logger    *zap.Logger
}

// NewEventService 创建事件服务
func NewEventService(publisher publisher.EventPublisher, logger *zap.Logger) *EventService {
	return &EventService{
		publisher: publisher,
		logger:    logger,
	}
}

// PublishUserRegistered 发布用户注册事件
func (s *EventService) PublishUserRegistered(ctx context.Context, userID string, username, email string) error {
	event := &types.UserRegisteredEvent{
		BaseEvent: types.NewBaseEvent(
			types.UserRegistered,
			"user-center",
			s.getRequestID(ctx),
			userID,
		),
		Username: username,
		Email:    email,
	}

	s.logger.Info("Publishing user registered event",
		zap.String("user_id", userID),
		zap.String("username", username),
		zap.String("email", email),
	)

	return s.publisher.PublishUserRegistered(ctx, event)
}

// PublishUserLoggedIn 发布用户登录事件
func (s *EventService) PublishUserLoggedIn(ctx context.Context, userID, username, ipAddress string) error {
	event := &types.UserLoggedInEvent{
		BaseEvent: types.NewBaseEvent(
			types.UserLoggedIn,
			"user-center",
			s.getRequestID(ctx),
			userID,
		),
		Username:  username,
		IPAddress: ipAddress,
	}

	s.logger.Info("Publishing user logged in event",
		zap.String("user_id", userID),
		zap.String("username", username),
		zap.String("ip_address", ipAddress),
	)

	return s.publisher.PublishUserLoggedIn(ctx, event)
}

// PublishUserPasswordChanged 发布用户密码变更事件
func (s *EventService) PublishUserPasswordChanged(ctx context.Context, userID string) error {
	event := &types.UserPasswordChangedEvent{
		BaseEvent: types.NewBaseEvent(
			types.UserPasswordChanged,
			"user-center",
			s.getRequestID(ctx),
			userID,
		),
	}

	s.logger.Info("Publishing user password changed event",
		zap.String("user_id", userID),
	)

	return s.publisher.PublishUserPasswordChanged(ctx, event)
}

// PublishUserStatusChanged 发布用户状态变更事件
func (s *EventService) PublishUserStatusChanged(ctx context.Context, userID, oldStatus, newStatus string) error {
	event := &types.UserStatusChangedEvent{
		BaseEvent: types.NewBaseEvent(
			types.UserStatusChanged,
			"user-center",
			s.getRequestID(ctx),
			userID,
		),
		OldStatus: oldStatus,
		NewStatus: newStatus,
	}

	s.logger.Info("Publishing user status changed event",
		zap.String("user_id", userID),
		zap.String("old_status", oldStatus),
		zap.String("new_status", newStatus),
	)

	return s.publisher.PublishUserStatusChanged(ctx, event)
}

// PublishUserDeleted 发布用户删除事件
func (s *EventService) PublishUserDeleted(ctx context.Context, userID string) error {
	event := &types.UserDeletedEvent{
		BaseEvent: types.NewBaseEvent(
			types.UserDeleted,
			"user-center",
			s.getRequestID(ctx),
			userID,
		),
	}

	s.logger.Info("Publishing user deleted event",
		zap.String("user_id", userID),
	)

	return s.publisher.PublishUserDeleted(ctx, event)
}

// PublishUserUpdated 发布用户更新事件
func (s *EventService) PublishUserUpdated(ctx context.Context, userID string) error {
	event := &types.UserUpdatedEvent{
		BaseEvent: types.NewBaseEvent(
			types.UserUpdated,
			"user-center",
			s.getRequestID(ctx),
			userID,
		),
		Changes: make(map[string]interface{}),
	}

	s.logger.Info("Publishing user updated event",
		zap.String("user_id", userID),
	)

	return s.publisher.PublishUserUpdated(ctx, event)
}

// getRequestID 获取请求ID
func (s *EventService) getRequestID(ctx context.Context) string {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		if requestID, exists := ginCtx.Get("request_id"); exists {
			if id, ok := requestID.(string); ok {
				return id
			}
		}
	}
	return ""
}
