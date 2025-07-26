package service

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/zhwjimmy/user-center/internal/infrastructure/messaging"
	"github.com/zhwjimmy/user-center/internal/kafka/event"
	"github.com/zhwjimmy/user-center/internal/model"
	"go.uber.org/zap"
)

// EventService provides event publishing services
type EventService struct {
	kafkaService messaging.Service
	logger       *zap.Logger
}

// NewEventService creates a new event service
func NewEventService(kafkaService messaging.Service, logger *zap.Logger) *EventService {
	return &EventService{
		kafkaService: kafkaService,
		logger:       logger,
	}
}

// PublishUserRegisteredEvent publishes a user registered event
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

// PublishUserLoggedInEvent publishes a user logged in event
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

// PublishUserPasswordChangedEvent publishes a user password changed event
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

// PublishUserStatusChangedEvent publishes a user status changed event
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

// PublishUserDeletedEvent publishes a user deleted event
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

// PublishUserUpdatedEvent publishes a user updated event
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

// getRequestID extracts request ID from context
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

// getStringValue safely gets string value from pointer
func (s *EventService) getStringValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}
