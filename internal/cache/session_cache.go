package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/zhwjimmy/user-center/internal/infrastructure/cache"
)

// SessionCache 会话缓存服务
type SessionCache struct {
	cache cache.Cache
}

// NewSessionCache 创建会话缓存服务
func NewSessionCache(cache cache.Cache) *SessionCache {
	return &SessionCache{
		cache: cache,
	}
}

// CacheSession 缓存会话数据
func (sc *SessionCache) CacheSession(ctx context.Context, sessionID string, session interface{}, expiration time.Duration) error {
	key := fmt.Sprintf("%s%s", SessionCacheKeyPrefix, sessionID)
	return sc.cache.Set(ctx, key, session, expiration)
}

// GetCachedSession 获取缓存的会话数据
func (sc *SessionCache) GetCachedSession(ctx context.Context, sessionID string, dest interface{}) error {
	key := fmt.Sprintf("%s%s", SessionCacheKeyPrefix, sessionID)
	return sc.cache.Get(ctx, key, dest)
}

// InvalidateSession 清除会话缓存
func (sc *SessionCache) InvalidateSession(ctx context.Context, sessionID string) error {
	key := fmt.Sprintf("%s%s", SessionCacheKeyPrefix, sessionID)
	return sc.cache.Delete(ctx, key)
}

const (
	SessionCacheKeyPrefix = "session:"
)
