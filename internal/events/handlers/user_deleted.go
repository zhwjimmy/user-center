package handlers

import (
	"context"
	"encoding/json"

	"github.com/zhwjimmy/user-center/internal/events/types"
	"go.uber.org/zap"
)

// UserDeletedHandler 用户删除事件处理器
type UserDeletedHandler struct {
	logger *zap.Logger
}

// NewUserDeletedHandler 创建用户删除事件处理器
func NewUserDeletedHandler(logger *zap.Logger) *UserDeletedHandler {
	return &UserDeletedHandler{
		logger: logger,
	}
}

// Handle 处理用户删除事件
func (h *UserDeletedHandler) Handle(ctx context.Context, payload []byte, headers map[string][]byte) error {
	var event types.UserDeletedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	h.logger.Info("Processing user deleted event",
		zap.String("user_id", event.UserID),
		zap.String("username", event.Username),
		zap.String("request_id", event.RequestID),
	)

	// 业务逻辑处理
	// 1. 清理用户数据
	if err := h.cleanupUserData(ctx, &event); err != nil {
		h.logger.Error("Failed to cleanup user data",
			zap.String("user_id", event.UserID),
			zap.Error(err),
		)
	}

	// 2. 发送账户删除确认
	if err := h.sendAccountDeletionConfirmation(ctx, &event); err != nil {
		h.logger.Error("Failed to send account deletion confirmation",
			zap.String("user_id", event.UserID),
			zap.Error(err),
		)
	}

	return nil
}

// cleanupUserData 清理用户数据
func (h *UserDeletedHandler) cleanupUserData(ctx context.Context, event *types.UserDeletedEvent) error {
	// TODO: 实现清理用户数据逻辑
	h.logger.Info("Cleaning up user data", zap.String("user_id", event.UserID))
	return nil
}

// sendAccountDeletionConfirmation 发送账户删除确认
func (h *UserDeletedHandler) sendAccountDeletionConfirmation(ctx context.Context, event *types.UserDeletedEvent) error {
	// TODO: 实现发送账户删除确认逻辑
	h.logger.Info("Sending account deletion confirmation", zap.String("user_id", event.UserID))
	return nil
}
