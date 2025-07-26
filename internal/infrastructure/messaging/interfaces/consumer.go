package interfaces

import (
	"context"
)

// Consumer 消费者接口
type Consumer interface {
	Start(ctx context.Context) error
	Stop() error
}
