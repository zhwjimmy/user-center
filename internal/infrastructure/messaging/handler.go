package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/zhwjimmy/user-center/internal/kafka/event"
	"go.uber.org/zap"
)

// UserEventHandler 用户事件处理器
type UserEventHandler struct {
	logger *zap.Logger
	// 可以注入其他服务，如邮件服务、通知服务等
}

// NewUserEventHandler 创建用户事件处理器
func NewUserEventHandler(logger *zap.Logger) MessageHandler {
	return &UserEventHandler{
		logger: logger,
	}
}

// HandleMessage 处理消息
func (h *UserEventHandler) HandleMessage(ctx context.Context, message *sarama.ConsumerMessage) error {
	// 获取事件类型
	eventType := h.getEventType(message.Headers)

	h.logger.Debug("Processing message",
		zap.String("topic", message.Topic),
		zap.String("event_type", eventType),
		zap.Int32("partition", message.Partition),
		zap.Int64("offset", message.Offset),
	)

	switch event.EventType(eventType) {
	case event.UserRegistered:
		var userEvent event.UserRegisteredEvent
		if err := json.Unmarshal(message.Value, &userEvent); err != nil {
			return fmt.Errorf("failed to unmarshal user registered event: %w", err)
		}
		return h.HandleUserRegistered(ctx, &userEvent)

	case event.UserLoggedIn:
		var userEvent event.UserLoggedInEvent
		if err := json.Unmarshal(message.Value, &userEvent); err != nil {
			return fmt.Errorf("failed to unmarshal user logged in event: %w", err)
		}
		return h.HandleUserLoggedIn(ctx, &userEvent)

	case event.UserPasswordChanged:
		var userEvent event.UserPasswordChangedEvent
		if err := json.Unmarshal(message.Value, &userEvent); err != nil {
			return fmt.Errorf("failed to unmarshal user password changed event: %w", err)
		}
		return h.HandleUserPasswordChanged(ctx, &userEvent)

	case event.UserStatusChanged:
		var userEvent event.UserStatusChangedEvent
		if err := json.Unmarshal(message.Value, &userEvent); err != nil {
			return fmt.Errorf("failed to unmarshal user status changed event: %w", err)
		}
		return h.HandleUserStatusChanged(ctx, &userEvent)

	case event.UserDeleted:
		var userEvent event.UserDeletedEvent
		if err := json.Unmarshal(message.Value, &userEvent); err != nil {
			return fmt.Errorf("failed to unmarshal user deleted event: %w", err)
		}
		return h.HandleUserDeleted(ctx, &userEvent)

	case event.UserUpdated:
		var userEvent event.UserUpdatedEvent
		if err := json.Unmarshal(message.Value, &userEvent); err != nil {
			return fmt.Errorf("failed to unmarshal user updated event: %w", err)
		}
		return h.HandleUserUpdated(ctx, &userEvent)

	default:
		h.logger.Warn("Unknown event type", zap.String("event_type", eventType))
		return nil // 忽略未知事件类型
	}
}

// getEventType 从消息头获取事件类型
func (h *UserEventHandler) getEventType(headers []*sarama.RecordHeader) string {
	for _, header := range headers {
		if string(header.Key) == "event_type" {
			return string(header.Value)
		}
	}
	return ""
}

// HandleUserRegistered 处理用户注册事件
func (h *UserEventHandler) HandleUserRegistered(ctx context.Context, event *event.UserRegisteredEvent) error {
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

// HandleUserLoggedIn 处理用户登录事件
func (h *UserEventHandler) HandleUserLoggedIn(ctx context.Context, event *event.UserLoggedInEvent) error {
	h.logger.Info("Processing user logged in event",
		zap.String("user_id", event.UserID),
		zap.String("username", event.Username),
		zap.String("ip_address", event.IPAddress),
		zap.String("request_id", event.RequestID),
	)

	// 业务逻辑处理
	// 1. 记录登录日志
	if err := h.recordLoginLog(ctx, event); err != nil {
		h.logger.Error("Failed to record login log",
			zap.String("user_id", event.UserID),
			zap.Error(err),
		)
	}

	// 2. 更新最后登录时间
	if err := h.updateLastLoginTime(ctx, event); err != nil {
		h.logger.Error("Failed to update last login time",
			zap.String("user_id", event.UserID),
			zap.Error(err),
		)
	}

	// 3. 检查异常登录
	if err := h.checkAnomalousLogin(ctx, event); err != nil {
		h.logger.Error("Failed to check anomalous login",
			zap.String("user_id", event.UserID),
			zap.Error(err),
		)
	}

	return nil
}

// HandleUserPasswordChanged 处理用户密码变更事件
func (h *UserEventHandler) HandleUserPasswordChanged(ctx context.Context, event *event.UserPasswordChangedEvent) error {
	h.logger.Info("Processing user password changed event",
		zap.String("user_id", event.UserID),
		zap.String("request_id", event.RequestID),
	)

	// 业务逻辑处理
	// 1. 发送密码变更通知
	if err := h.sendPasswordChangeNotification(ctx, event); err != nil {
		h.logger.Error("Failed to send password change notification",
			zap.String("user_id", event.UserID),
			zap.Error(err),
		)
	}

	// 2. 记录安全日志
	if err := h.recordSecurityLog(ctx, event); err != nil {
		h.logger.Error("Failed to record security log",
			zap.String("user_id", event.UserID),
			zap.Error(err),
		)
	}

	return nil
}

// HandleUserStatusChanged 处理用户状态变更事件
func (h *UserEventHandler) HandleUserStatusChanged(ctx context.Context, event *event.UserStatusChangedEvent) error {
	h.logger.Info("Processing user status changed event",
		zap.String("user_id", event.UserID),
		zap.String("old_status", event.OldStatus),
		zap.String("new_status", event.NewStatus),
		zap.String("request_id", event.RequestID),
	)

	// 业务逻辑处理
	// 1. 发送状态变更通知
	if err := h.sendStatusChangeNotification(ctx, event); err != nil {
		h.logger.Error("Failed to send status change notification",
			zap.String("user_id", event.UserID),
			zap.Error(err),
		)
	}

	// 2. 更新用户状态缓存
	if err := h.updateUserStatusCache(ctx, event); err != nil {
		h.logger.Error("Failed to update user status cache",
			zap.String("user_id", event.UserID),
			zap.Error(err),
		)
	}

	return nil
}

// HandleUserDeleted 处理用户删除事件
func (h *UserEventHandler) HandleUserDeleted(ctx context.Context, event *event.UserDeletedEvent) error {
	h.logger.Info("Processing user deleted event",
		zap.String("user_id", event.UserID),
		zap.String("request_id", event.RequestID),
	)

	// 业务逻辑处理
	// 1. 清理用户数据
	if err := h.cleanupUserData(ctx, event); err != nil {
		h.logger.Error("Failed to cleanup user data",
			zap.String("user_id", event.UserID),
			zap.Error(err),
		)
	}

	// 2. 发送账户删除确认
	if err := h.sendAccountDeletionConfirmation(ctx, event); err != nil {
		h.logger.Error("Failed to send account deletion confirmation",
			zap.String("user_id", event.UserID),
			zap.Error(err),
		)
	}

	return nil
}

// HandleUserUpdated 处理用户更新事件
func (h *UserEventHandler) HandleUserUpdated(ctx context.Context, event *event.UserUpdatedEvent) error {
	h.logger.Info("Processing user updated event",
		zap.String("user_id", event.UserID),
		zap.String("request_id", event.RequestID),
	)

	// 业务逻辑处理
	// 1. 更新用户缓存
	if err := h.updateUserCache(ctx, event); err != nil {
		h.logger.Error("Failed to update user cache",
			zap.String("user_id", event.UserID),
			zap.Error(err),
		)
	}

	// 2. 同步用户信息到外部系统
	if err := h.syncUserInfoToExternalSystems(ctx, event); err != nil {
		h.logger.Error("Failed to sync user info to external systems",
			zap.String("user_id", event.UserID),
			zap.Error(err),
		)
	}

	return nil
}

// 以下是具体的业务逻辑实现方法（简化版本）

func (h *UserEventHandler) sendWelcomeEmail(ctx context.Context, event *event.UserRegisteredEvent) error {
	// TODO: 实现发送欢迎邮件逻辑
	return nil
}

func (h *UserEventHandler) initializeUserSettings(ctx context.Context, event *event.UserRegisteredEvent) error {
	// TODO: 实现初始化用户设置逻辑
	return nil
}

func (h *UserEventHandler) recordUserRegistrationStats(ctx context.Context, event *event.UserRegisteredEvent) error {
	// TODO: 实现记录用户注册统计逻辑
	return nil
}

func (h *UserEventHandler) recordLoginLog(ctx context.Context, event *event.UserLoggedInEvent) error {
	// TODO: 实现记录登录日志逻辑
	return nil
}

func (h *UserEventHandler) updateLastLoginTime(ctx context.Context, event *event.UserLoggedInEvent) error {
	// TODO: 实现更新最后登录时间逻辑
	return nil
}

func (h *UserEventHandler) checkAnomalousLogin(ctx context.Context, event *event.UserLoggedInEvent) error {
	// TODO: 实现检查异常登录逻辑
	return nil
}

func (h *UserEventHandler) sendPasswordChangeNotification(ctx context.Context, event *event.UserPasswordChangedEvent) error {
	// TODO: 实现发送密码变更通知逻辑
	return nil
}

func (h *UserEventHandler) recordSecurityLog(ctx context.Context, event *event.UserPasswordChangedEvent) error {
	// TODO: 实现记录安全日志逻辑
	return nil
}

func (h *UserEventHandler) sendStatusChangeNotification(ctx context.Context, event *event.UserStatusChangedEvent) error {
	// TODO: 实现发送状态变更通知逻辑
	return nil
}

func (h *UserEventHandler) updateUserStatusCache(ctx context.Context, event *event.UserStatusChangedEvent) error {
	// TODO: 实现更新用户状态缓存逻辑
	return nil
}

func (h *UserEventHandler) cleanupUserData(ctx context.Context, event *event.UserDeletedEvent) error {
	// TODO: 实现清理用户数据逻辑
	return nil
}

func (h *UserEventHandler) sendAccountDeletionConfirmation(ctx context.Context, event *event.UserDeletedEvent) error {
	// TODO: 实现发送账户删除确认逻辑
	return nil
}

func (h *UserEventHandler) updateUserCache(ctx context.Context, event *event.UserUpdatedEvent) error {
	// TODO: 实现更新用户缓存逻辑
	return nil
}

func (h *UserEventHandler) syncUserInfoToExternalSystems(ctx context.Context, event *event.UserUpdatedEvent) error {
	// TODO: 实现同步用户信息到外部系统逻辑
	return nil
}
