package handlers

import (
	"context"

	"github.com/zhwjimmy/user-center/internal/events/types"
	"go.uber.org/zap"
)

// UserRegisteredHandler 用户注册事件处理器
type UserRegisteredHandler struct {
	logger *zap.Logger
	// 可以注入其他服务，如邮件服务、通知服务等
}

// NewUserRegisteredHandler 创建用户注册事件处理器
func NewUserRegisteredHandler(logger *zap.Logger) *UserRegisteredHandler {
	return &UserRegisteredHandler{
		logger: logger,
	}
}

// Handle 处理用户注册事件
func (h *UserRegisteredHandler) Handle(ctx context.Context, event *types.UserRegisteredEvent) error {
	h.logger.Info("Processing user registered event",
		zap.String("user_id", event.UserID),
		zap.String("username", event.Username),
		zap.String("email", event.Email),
		zap.String("request_id", event.RequestID),
	)

	// 业务逻辑处理
	// 1. 发送欢迎邮件
	if err := h.sendWelcomeEmail(ctx, event); err != nil {
		h.logger.Error("Failed to send welcome email",
			zap.String("user_id", event.UserID),
			zap.Error(err),
		)
		// 不返回错误，避免阻塞消息处理
	}

	// 2. 初始化用户配置
	if err := h.initializeUserSettings(ctx, event); err != nil {
		h.logger.Error("Failed to initialize user settings",
			zap.String("user_id", event.UserID),
			zap.Error(err),
		)
	}

	// 3. 记录用户注册统计
	if err := h.recordUserRegistrationStats(ctx, event); err != nil {
		h.logger.Error("Failed to record user registration stats",
			zap.String("user_id", event.UserID),
			zap.Error(err),
		)
	}

	return nil
}

// sendWelcomeEmail 发送欢迎邮件
func (h *UserRegisteredHandler) sendWelcomeEmail(ctx context.Context, event *types.UserRegisteredEvent) error {
	// TODO: 实现发送欢迎邮件逻辑
	h.logger.Info("Sending welcome email", zap.String("email", event.Email))
	return nil
}

// initializeUserSettings 初始化用户设置
func (h *UserRegisteredHandler) initializeUserSettings(ctx context.Context, event *types.UserRegisteredEvent) error {
	// TODO: 实现初始化用户设置逻辑
	h.logger.Info("Initializing user settings", zap.String("user_id", event.UserID))
	return nil
}

// recordUserRegistrationStats 记录用户注册统计
func (h *UserRegisteredHandler) recordUserRegistrationStats(ctx context.Context, event *types.UserRegisteredEvent) error {
	// TODO: 实现记录用户注册统计逻辑
	h.logger.Info("Recording user registration stats", zap.String("user_id", event.UserID))
	return nil
}
