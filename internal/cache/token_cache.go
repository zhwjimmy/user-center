package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/zhwjimmy/user-center/internal/infrastructure/cache"
)

// TokenCache 令牌缓存服务
type TokenCache struct {
	cache cache.Cache
}

// NewTokenCache 创建令牌缓存服务
func NewTokenCache(cache cache.Cache) *TokenCache {
	return &TokenCache{
		cache: cache,
	}
}

// BlacklistToken 将令牌加入黑名单
func (tc *TokenCache) BlacklistToken(ctx context.Context, token string, expiration time.Duration) error {
	key := fmt.Sprintf("%s%s", TokenBlacklistPrefix, token)
	return tc.cache.Set(ctx, key, true, expiration)
}

// IsTokenBlacklisted 检查令牌是否在黑名单中
func (tc *TokenCache) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	key := fmt.Sprintf("%s%s", TokenBlacklistPrefix, token)
	return tc.cache.Exists(ctx, key)
}

// RemoveFromBlacklist 从黑名单中移除令牌
func (tc *TokenCache) RemoveFromBlacklist(ctx context.Context, token string) error {
	key := fmt.Sprintf("%s%s", TokenBlacklistPrefix, token)
	return tc.cache.Delete(ctx, key)
}

const (
	TokenBlacklistPrefix = "token_blacklist:"
)
