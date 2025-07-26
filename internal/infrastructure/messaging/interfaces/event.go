package interfaces

// Event 通用事件接口
type Event interface {
	GetTopic() string
	GetEventType() string
	GetUserID() string
	GetRequestID() string
	GetTimestamp() string
}
