package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/zhwjimmy/user-center/internal/infrastructure/cache"
)

// UserCache 用户缓存服务
type UserCache struct {
	cache cache.Cache
}

// NewUserCache 创建用户缓存服务
func NewUserCache(cache cache.Cache) *UserCache {
	return &UserCache{
		cache: cache,
	}
}

// CacheUser 缓存用户数据
func (uc *UserCache) CacheUser(ctx context.Context, userID uint, user interface{}, expiration time.Duration) error {
	key := fmt.Sprintf("%s%d", UserCacheKeyPrefix, userID)
	return uc.cache.Set(ctx, key, user, expiration)
}

// GetCachedUser 获取缓存的用户数据
func (uc *UserCache) GetCachedUser(ctx context.Context, userID uint, dest interface{}) error {
	key := fmt.Sprintf("%s%d", UserCacheKeyPrefix, userID)
	return uc.cache.Get(ctx, key, dest)
}

// InvalidateUserCache 清除用户缓存
func (uc *UserCache) InvalidateUserCache(ctx context.Context, userID uint) error {
	key := fmt.Sprintf("%s%d", UserCacheKeyPrefix, userID)
	return uc.cache.Delete(ctx, key)
}

// CacheKey 常量
const (
	UserCacheKeyPrefix = "user:"
)
