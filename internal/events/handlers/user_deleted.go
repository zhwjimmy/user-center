package handlers

import (
	"context"

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
func (h *UserDeletedHandler) Handle(ctx context.Context, event *types.UserDeletedEvent) error {
	h.logger.Info("Processing user deleted event",
		zap.String("user_id", event.UserID),
		zap.String("request_id", event.RequestID),
	)

	// TODO: 实现用户删除相关业务逻辑
	// 1. 清理用户数据
	// 2. 发送账户删除确认
	// 3. 记录删除日志

	return nil
}
