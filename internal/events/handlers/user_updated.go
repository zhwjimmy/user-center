package handlers

import (
	"context"
	"encoding/json"

	"github.com/zhwjimmy/user-center/internal/events/types"
	"go.uber.org/zap"
)

// UserUpdatedHandler 用户更新事件处理器
type UserUpdatedHandler struct {
	logger *zap.Logger
}

// NewUserUpdatedHandler 创建用户更新事件处理器
func NewUserUpdatedHandler(logger *zap.Logger) *UserUpdatedHandler {
	return &UserUpdatedHandler{
		logger: logger,
	}
}

// Handle 处理用户更新事件
func (h *UserUpdatedHandler) Handle(ctx context.Context, payload []byte, headers map[string][]byte) error {
	var event types.UserUpdatedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	h.logger.Info("Processing user updated event",
		zap.String("user_id", event.UserID),
		zap.String("username", event.Username),
		zap.String("request_id", event.RequestID),
	)

	// 业务逻辑处理
	// 1. 更新用户缓存
	if err := h.updateUserCache(ctx, &event); err != nil {
		h.logger.Error("Failed to update user cache",
			zap.String("user_id", event.UserID),
			zap.Error(err),
		)
	}

	// 2. 同步用户信息到外部系统
	if err := h.syncUserInfoToExternalSystems(ctx, &event); err != nil {
		h.logger.Error("Failed to sync user info to external systems",
			zap.String("user_id", event.UserID),
			zap.Error(err),
		)
	}

	return nil
}

// updateUserCache 更新用户缓存
func (h *UserUpdatedHandler) updateUserCache(ctx context.Context, event *types.UserUpdatedEvent) error {
	// TODO: 实现更新用户缓存逻辑
	h.logger.Info("Updating user cache", zap.String("user_id", event.UserID))
	return nil
}

// syncUserInfoToExternalSystems 同步用户信息到外部系统
func (h *UserUpdatedHandler) syncUserInfoToExternalSystems(ctx context.Context, event *types.UserUpdatedEvent) error {
	// TODO: 实现同步用户信息到外部系统逻辑
	h.logger.Info("Syncing user info to external systems", zap.String("user_id", event.UserID))
	return nil
}
