package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/zhwjimmy/user-center/internal/infrastructure/cache"
)

// RateLimitCache 速率限制缓存服务
type RateLimitCache struct {
	cache cache.Cache
}

// NewRateLimitCache 创建速率限制缓存服务
func NewRateLimitCache(cache cache.Cache) *RateLimitCache {
	return &RateLimitCache{
		cache: cache,
	}
}

// SetRateLimit 设置速率限制计数器
func (rlc *RateLimitCache) SetRateLimit(ctx context.Context, identifier string, expiration time.Duration) (int64, error) {
	key := fmt.Sprintf("%s%s", RateLimitKeyPrefix, identifier)
	return rlc.cache.IncrementWithExpiry(ctx, key, expiration)
}

// CheckRateLimit 检查速率限制
func (rlc *RateLimitCache) CheckRateLimit(ctx context.Context, identifier string, limit int64, expiration time.Duration) (bool, error) {
	count, err := rlc.SetRateLimit(ctx, identifier, expiration)
	if err != nil {
		return false, err
	}
	return count <= limit, nil
}

const (
	RateLimitKeyPrefix = "rate_limit:"
)
