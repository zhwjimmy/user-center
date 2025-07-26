// Package cache provides caching interfaces and implementations for the user-center application.
// It defines a common Cache interface that can be implemented by different caching backends
// such as Redis, in-memory cache, etc.
package cache

import (
	"context"
	"time"
)

// Cache 缓存接口
type Cache interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string, dest interface{}) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error)
	Increment(ctx context.Context, key string) (int64, error)
	IncrementWithExpiry(ctx context.Context, key string, expiration time.Duration) (int64, error)
	SetExpiry(ctx context.Context, key string, expiration time.Duration) error
	GetTTL(ctx context.Context, key string) (time.Duration, error)
	Keys(ctx context.Context, pattern string) ([]string, error)
	Close() error
	Health(ctx context.Context) error
}
