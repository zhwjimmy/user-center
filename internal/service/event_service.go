package service

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/zhwjimmy/user-center/internal/kafka"
	"github.com/zhwjimmy/user-center/internal/kafka/event"
	"github.com/zhwjimmy/user-center/internal/model"
	"go.uber.org/zap"
)

// EventService 事件服务
type EventService struct {
	kafkaService kafka.Service
	logger       *zap.Logger
}

// NewEventService 创建事件服务
func NewEventService(kafkaService kafka.Service, logger *zap.Logger) *EventService {
	return &EventService{
		kafkaService: kafkaService,
		logger:       logger,
	}
}

// PublishUserRegisteredEvent 发布用户注册事件
func (s *EventService) PublishUserRegisteredEvent(ctx context.Context, user *model.User) error {
	requestID := s.getRequestID(ctx)

	userEvent := &event.UserRegisteredEvent{
		BaseEvent: event.NewBaseEvent(
			event.UserRegistered,
			"user-center",
			requestID,
			user.ID,
		),
		Username:  user.Username,
		Email:     user.Email,
		FirstName: s.getStringValue(user.FirstName),
		LastName:  s.getStringValue(user.LastName),
	}

	return s.kafkaService.GetProducer().PublishUserEventAsync(ctx, userEvent)
}

// PublishUserLoggedInEvent 发布用户登录事件
func (s *EventService) PublishUserLoggedInEvent(ctx context.Context, user *model.User, ipAddress, userAgent string) error {
	requestID := s.getRequestID(ctx)

	userEvent := &event.UserLoggedInEvent{
		BaseEvent: event.NewBaseEvent(
			event.UserLoggedIn,
			"user-center",
			requestID,
			user.ID,
		),
		Username:  user.Username,
		Email:     user.Email,
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}

	return s.kafkaService.GetProducer().PublishUserEventAsync(ctx, userEvent)
}

// PublishUserPasswordChangedEvent 发布用户密码变更事件
func (s *EventService) PublishUserPasswordChangedEvent(ctx context.Context, user *model.User, ipAddress string) error {
	requestID := s.getRequestID(ctx)

	userEvent := &event.UserPasswordChangedEvent{
		BaseEvent: event.NewBaseEvent(
			event.UserPasswordChanged,
			"user-center",
			requestID,
			user.ID,
		),
		Username:  user.Username,
		Email:     user.Email,
		IPAddress: ipAddress,
	}

	return s.kafkaService.GetProducer().PublishUserEventAsync(ctx, userEvent)
}

// PublishUserStatusChangedEvent 发布用户状态变更事件
func (s *EventService) PublishUserStatusChangedEvent(ctx context.Context, user *model.User, oldStatus, newStatus string) error {
	requestID := s.getRequestID(ctx)

	userEvent := &event.UserStatusChangedEvent{
		BaseEvent: event.NewBaseEvent(
			event.UserStatusChanged,
			"user-center",
			requestID,
			user.ID,
		),
		Username:  user.Username,
		Email:     user.Email,
		OldStatus: oldStatus,
		NewStatus: newStatus,
	}

	return s.kafkaService.GetProducer().PublishUserEventAsync(ctx, userEvent)
}

// PublishUserDeletedEvent 发布用户删除事件
func (s *EventService) PublishUserDeletedEvent(ctx context.Context, user *model.User) error {
	requestID := s.getRequestID(ctx)

	userEvent := &event.UserDeletedEvent{
		BaseEvent: event.NewBaseEvent(
			event.UserDeleted,
			"user-center",
			requestID,
			user.ID,
		),
		Username: user.Username,
		Email:    user.Email,
	}

	return s.kafkaService.GetProducer().PublishUserEventAsync(ctx, userEvent)
}

// PublishUserUpdatedEvent 发布用户更新事件
func (s *EventService) PublishUserUpdatedEvent(ctx context.Context, user *model.User, changes map[string]interface{}) error {
	requestID := s.getRequestID(ctx)

	userEvent := &event.UserUpdatedEvent{
		BaseEvent: event.NewBaseEvent(
			event.UserUpdated,
			"user-center",
			requestID,
			user.ID,
		),
		Username: user.Username,
		Email:    user.Email,
		Changes:  changes,
	}

	return s.kafkaService.GetProducer().PublishUserEventAsync(ctx, userEvent)
}

// getRequestID 从上下文获取请求ID
func (s *EventService) getRequestID(ctx context.Context) string {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		if requestID := ginCtx.GetHeader("X-Request-ID"); requestID != "" {
			return requestID
		}
	}
	return ""
}

// getStringValue 获取字符串指针的值
func (s *EventService) getStringValue(ptr *string) string {
	if ptr != nil {
		return *ptr
	}
	return ""
}
