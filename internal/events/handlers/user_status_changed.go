package handlers

import (
	"context"

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
func (h *UserStatusChangedHandler) Handle(ctx context.Context, event *types.UserStatusChangedEvent) error {
	h.logger.Info("Processing user status changed event",
		zap.String("user_id", event.UserID),
		zap.String("old_status", event.OldStatus),
		zap.String("new_status", event.NewStatus),
		zap.String("request_id", event.RequestID),
	)

	// TODO: 实现状态变更相关业务逻辑
	// 1. 发送状态变更通知
	// 2. 更新用户状态缓存
	// 3. 记录状态变更日志

	return nil
}
