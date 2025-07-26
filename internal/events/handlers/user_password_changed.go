package handlers

import (
	"context"

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
func (h *UserPasswordChangedHandler) Handle(ctx context.Context, event *types.UserPasswordChangedEvent) error {
	h.logger.Info("Processing user password changed event",
		zap.String("user_id", event.UserID),
		zap.String("request_id", event.RequestID),
	)

	// TODO: 实现密码变更相关业务逻辑
	// 1. 发送密码变更通知
	// 2. 记录安全日志
	// 3. 更新密码历史

	return nil
}
