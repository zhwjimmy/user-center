package interfaces

import (
	"context"
)

// Producer 生产者接口
type Producer interface {
	PublishUserEvent(ctx context.Context, event interface{}) error
	PublishUserEventAsync(ctx context.Context, event interface{}) error
	Close() error
}
