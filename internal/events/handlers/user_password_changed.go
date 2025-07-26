package handlers

import (
	"context"
	"encoding/json"

	"github.com/zhwjimmy/user-center/internal/events/types"
	"go.uber.org/zap"
)

// UserPasswordChangedHandler 用户密码变更事件处理器
type UserPasswordChangedHandler struct {
	logger *zap.Logger
}

// NewUserPasswordChangedHandler 创建用户密码变更事件处理器
func NewUserPasswordChangedHandler(logger *zap.Logger) *UserPasswordChangedHandler {
	return &UserPasswordChangedHandler{
		logger: logger,
	}
}

// Handle 处理用户密码变更事件
func (h *UserPasswordChangedHandler) Handle(ctx context.Context, payload []byte, headers map[string][]byte) error {
	var event types.UserPasswordChangedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	h.logger.Info("Processing user password changed event",
		zap.String("user_id", event.UserID),
		zap.String("username", event.Username),
		zap.String("request_id", event.RequestID),
	)

	// 业务逻辑处理
	// 1. 发送密码变更通知
	if err := h.sendPasswordChangeNotification(ctx, &event); err != nil {
		h.logger.Error("Failed to send password change notification",
			zap.String("user_id", event.UserID),
			zap.Error(err),
		)
	}

	// 2. 记录安全日志
	if err := h.recordSecurityLog(ctx, &event); err != nil {
		h.logger.Error("Failed to record security log",
			zap.String("user_id", event.UserID),
			zap.Error(err),
		)
	}

	return nil
}

// sendPasswordChangeNotification 发送密码变更通知
func (h *UserPasswordChangedHandler) sendPasswordChangeNotification(ctx context.Context, event *types.UserPasswordChangedEvent) error {
	// TODO: 实现发送密码变更通知逻辑
	h.logger.Info("Sending password change notification", zap.String("user_id", event.UserID))
	return nil
}

// recordSecurityLog 记录安全日志
func (h *UserPasswordChangedHandler) recordSecurityLog(ctx context.Context, event *types.UserPasswordChangedEvent) error {
	// TODO: 实现记录安全日志逻辑
	h.logger.Info("Recording security log", zap.String("user_id", event.UserID))
	return nil
}
