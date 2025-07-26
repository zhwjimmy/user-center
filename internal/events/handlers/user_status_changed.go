package handlers

import (
	"context"
	"encoding/json"

	"github.com/zhwjimmy/user-center/internal/events/types"
	"go.uber.org/zap"
)

// UserStatusChangedHandler 用户状态变更事件处理器
type UserStatusChangedHandler struct {
	logger *zap.Logger
}

// NewUserStatusChangedHandler 创建用户状态变更事件处理器
func NewUserStatusChangedHandler(logger *zap.Logger) *UserStatusChangedHandler {
	return &UserStatusChangedHandler{
		logger: logger,
	}
}

// Handle 处理用户状态变更事件
func (h *UserStatusChangedHandler) Handle(ctx context.Context, payload []byte, headers map[string][]byte) error {
	var event types.UserStatusChangedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	h.logger.Info("Processing user status changed event",
		zap.String("user_id", event.UserID),
		zap.String("username", event.Username),
		zap.String("old_status", string(event.OldStatus)),
		zap.String("new_status", string(event.NewStatus)),
		zap.String("request_id", event.RequestID),
	)

	// 业务逻辑处理
	// 1. 发送状态变更通知
	if err := h.sendStatusChangeNotification(ctx, &event); err != nil {
		h.logger.Error("Failed to send status change notification",
			zap.String("user_id", event.UserID),
			zap.Error(err),
		)
	}

	// 2. 更新用户状态缓存
	if err := h.updateUserStatusCache(ctx, &event); err != nil {
		h.logger.Error("Failed to update user status cache",
			zap.String("user_id", event.UserID),
			zap.Error(err),
		)
	}

	return nil
}

// sendStatusChangeNotification 发送状态变更通知
func (h *UserStatusChangedHandler) sendStatusChangeNotification(ctx context.Context, event *types.UserStatusChangedEvent) error {
	// TODO: 实现发送状态变更通知逻辑
	h.logger.Info("Sending status change notification", zap.String("user_id", event.UserID))
	return nil
}

// updateUserStatusCache 更新用户状态缓存
func (h *UserStatusChangedHandler) updateUserStatusCache(ctx context.Context, event *types.UserStatusChangedEvent) error {
	// TODO: 实现更新用户状态缓存逻辑
	h.logger.Info("Updating user status cache", zap.String("user_id", event.UserID))
	return nil
}
