package event

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// EventType 定义事件类型
type EventType string

const (
	// 用户事件类型
	UserRegistered      EventType = "user.registered"
	UserLoggedIn        EventType = "user.logged_in"
	UserPasswordChanged EventType = "user.password_changed"
	UserStatusChanged   EventType = "user.status_changed"
	UserDeleted         EventType = "user.deleted"
	UserUpdated         EventType = "user.updated"
)

// BaseEvent 基础事件结构
type BaseEvent struct {
	ID        string                 `json:"id"`
	Type      EventType              `json:"type"`
	Source    string                 `json:"source"`
	Timestamp time.Time              `json:"timestamp"`
	Version   string                 `json:"version"`
	RequestID string                 `json:"request_id,omitempty"`
	UserID    string                 `json:"user_id,omitempty"`
	Data      map[string]interface{} `json:"data"`
}

// UserRegisteredEvent 用户注册事件
type UserRegisteredEvent struct {
	BaseEvent
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

// UserLoggedInEvent 用户登录事件
type UserLoggedInEvent struct {
	BaseEvent
	Username  string `json:"username"`
	Email     string `json:"email"`
	IPAddress string `json:"ip_address,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
}

// UserPasswordChangedEvent 用户密码变更事件
type UserPasswordChangedEvent struct {
	BaseEvent
	Username  string `json:"username"`
	Email     string `json:"email"`
	IPAddress string `json:"ip_address,omitempty"`
}

// UserStatusChangedEvent 用户状态变更事件
type UserStatusChangedEvent struct {
	BaseEvent
	Username  string `json:"username"`
	Email     string `json:"email"`
	OldStatus string `json:"old_status"`
	NewStatus string `json:"new_status"`
}

// UserDeletedEvent 用户删除事件
type UserDeletedEvent struct {
	BaseEvent
	Username string `json:"username"`
	Email    string `json:"email"`
}

// UserUpdatedEvent 用户更新事件
type UserUpdatedEvent struct {
	BaseEvent
	Username string                 `json:"username"`
	Email    string                 `json:"email"`
	Changes  map[string]interface{} `json:"changes"`
}

// NewBaseEvent 创建基础事件
func NewBaseEvent(eventType EventType, source, requestID, userID string) BaseEvent {
	return BaseEvent{
		ID:        generateEventID(),
		Type:      eventType,
		Source:    source,
		Timestamp: time.Now(),
		Version:   "1.0",
		RequestID: requestID,
		UserID:    userID,
		Data:      make(map[string]interface{}),
	}
}

// ToJSON 将事件转换为JSON
func (e *BaseEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// FromJSON 从JSON创建事件
func (e *BaseEvent) FromJSON(data []byte) error {
	return json.Unmarshal(data, e)
}

// generateEventID 生成事件ID
func generateEventID() string {
	return uuid.New().String()
}
