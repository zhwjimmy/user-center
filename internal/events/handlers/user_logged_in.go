package handlers

import (
	"context"
	"encoding/json"

	"github.com/zhwjimmy/user-center/internal/events/types"
	"go.uber.org/zap"
)

// UserLoggedInHandler 用户登录事件处理器
type UserLoggedInHandler struct {
	logger *zap.Logger
	// 可以注入其他服务，如安全服务、审计服务等
}

// NewUserLoggedInHandler 创建用户登录事件处理器
func NewUserLoggedInHandler(logger *zap.Logger) *UserLoggedInHandler {
	return &UserLoggedInHandler{
		logger: logger,
	}
}

// Handle 处理用户登录事件
func (h *UserLoggedInHandler) Handle(ctx context.Context, payload []byte, headers map[string][]byte) error {
	var event types.UserLoggedInEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	h.logger.Info("Processing user logged in event",
		zap.String("user_id", event.UserID),
		zap.String("username", event.Username),
		zap.String("ip_address", event.IPAddress),
		zap.String("request_id", event.RequestID),
	)

	// 业务逻辑处理
	// 1. 记录登录日志
	if err := h.recordLoginLog(ctx, &event); err != nil {
		h.logger.Error("Failed to record login log",
			zap.String("user_id", event.UserID),
			zap.Error(err),
		)
	}

	// 2. 更新最后登录时间
	if err := h.updateLastLoginTime(ctx, &event); err != nil {
		h.logger.Error("Failed to update last login time",
			zap.String("user_id", event.UserID),
			zap.Error(err),
		)
	}

	// 3. 检查异常登录
	if err := h.checkAnomalousLogin(ctx, &event); err != nil {
		h.logger.Error("Failed to check anomalous login",
			zap.String("user_id", event.UserID),
			zap.Error(err),
		)
	}

	return nil
}

// recordLoginLog 记录登录日志
func (h *UserLoggedInHandler) recordLoginLog(ctx context.Context, event *types.UserLoggedInEvent) error {
	// TODO: 实现记录登录日志逻辑
	h.logger.Info("Recording login log",
		zap.String("user_id", event.UserID),
		zap.String("ip_address", event.IPAddress),
	)
	return nil
}

// updateLastLoginTime 更新最后登录时间
func (h *UserLoggedInHandler) updateLastLoginTime(ctx context.Context, event *types.UserLoggedInEvent) error {
	// TODO: 实现更新最后登录时间逻辑
	h.logger.Info("Updating last login time", zap.String("user_id", event.UserID))
	return nil
}

// checkAnomalousLogin 检查异常登录
func (h *UserLoggedInHandler) checkAnomalousLogin(ctx context.Context, event *types.UserLoggedInEvent) error {
	// TODO: 实现检查异常登录逻辑
	h.logger.Info("Checking anomalous login",
		zap.String("user_id", event.UserID),
		zap.String("ip_address", event.IPAddress),
	)
	return nil
}
