package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/your-org/user-center/internal/config"
	"go.uber.org/zap"
)

// Redis represents Redis cache connection
type Redis struct {
	Client *redis.Client
	logger *zap.Logger
}

// NewRedis creates a new Redis connection
func NewRedis(cfg *config.Config, logger *zap.Logger) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.Addr,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.Info("Redis connected successfully",
		zap.String("addr", cfg.Redis.Addr),
		zap.Int("db", cfg.Redis.DB),
	)

	return &Redis{
		Client: client,
		logger: logger,
	}, nil
}

// Close closes the Redis connection
func (r *Redis) Close() error {
	return r.Client.Close()
}

// Health checks the Redis health
func (r *Redis) Health(ctx context.Context) error {
	return r.Client.Ping(ctx).Err()
}

// Set stores a value with expiration
func (r *Redis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	if err := r.Client.Set(ctx, key, data, expiration).Err(); err != nil {
		r.logger.Error("Failed to set cache",
			zap.String("key", key),
			zap.Error(err),
		)
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

// Get retrieves a value from cache
func (r *Redis) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("key not found")
		}
		r.logger.Error("Failed to get cache",
			zap.String("key", key),
			zap.Error(err),
		)
		return fmt.Errorf("failed to get cache: %w", err)
	}

	if err := json.Unmarshal([]byte(data), dest); err != nil {
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}

	return nil
}

// Delete removes a key from cache
func (r *Redis) Delete(ctx context.Context, key string) error {
	if err := r.Client.Del(ctx, key).Err(); err != nil {
		r.logger.Error("Failed to delete cache",
			zap.String("key", key),
			zap.Error(err),
		)
		return fmt.Errorf("failed to delete cache: %w", err)
	}

	return nil
}

// Exists checks if a key exists
func (r *Redis) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.Client.Exists(ctx, key).Result()
	if err != nil {
		r.logger.Error("Failed to check cache existence",
			zap.String("key", key),
			zap.Error(err),
		)
		return false, fmt.Errorf("failed to check cache existence: %w", err)
	}

	return result > 0, nil
}

// SetNX sets a value only if the key does not exist
func (r *Redis) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("failed to marshal value: %w", err)
	}

	result, err := r.Client.SetNX(ctx, key, data, expiration).Result()
	if err != nil {
		r.logger.Error("Failed to set cache NX",
			zap.String("key", key),
			zap.Error(err),
		)
		return false, fmt.Errorf("failed to set cache NX: %w", err)
	}

	return result, nil
}

// Increment increments a counter
func (r *Redis) Increment(ctx context.Context, key string) (int64, error) {
	result, err := r.Client.Incr(ctx, key).Result()
	if err != nil {
		r.logger.Error("Failed to increment counter",
			zap.String("key", key),
			zap.Error(err),
		)
		return 0, fmt.Errorf("failed to increment counter: %w", err)
	}

	return result, nil
}

// IncrementWithExpiry increments a counter with expiration
func (r *Redis) IncrementWithExpiry(ctx context.Context, key string, expiration time.Duration) (int64, error) {
	pipe := r.Client.Pipeline()
	incrCmd := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, expiration)

	if _, err := pipe.Exec(ctx); err != nil {
		r.logger.Error("Failed to increment counter with expiry",
			zap.String("key", key),
			zap.Error(err),
		)
		return 0, fmt.Errorf("failed to increment counter with expiry: %w", err)
	}

	return incrCmd.Val(), nil
}

// SetExpiry sets expiration for a key
func (r *Redis) SetExpiry(ctx context.Context, key string, expiration time.Duration) error {
	if err := r.Client.Expire(ctx, key, expiration).Err(); err != nil {
		r.logger.Error("Failed to set expiry",
			zap.String("key", key),
			zap.Error(err),
		)
		return fmt.Errorf("failed to set expiry: %w", err)
	}

	return nil
}

// GetTTL gets the remaining time to live of a key
func (r *Redis) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := r.Client.TTL(ctx, key).Result()
	if err != nil {
		r.logger.Error("Failed to get TTL",
			zap.String("key", key),
			zap.Error(err),
		)
		return 0, fmt.Errorf("failed to get TTL: %w", err)
	}

	return ttl, nil
}

// Keys returns all keys matching pattern
func (r *Redis) Keys(ctx context.Context, pattern string) ([]string, error) {
	keys, err := r.Client.Keys(ctx, pattern).Result()
	if err != nil {
		r.logger.Error("Failed to get keys",
			zap.String("pattern", pattern),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get keys: %w", err)
	}

	return keys, nil
}

// Cache key constants
const (
	UserCacheKeyPrefix    = "user:"
	SessionCacheKeyPrefix = "session:"
	RateLimitKeyPrefix    = "rate_limit:"
	TokenBlacklistPrefix  = "token_blacklist:"
)

// Helper functions for common cache operations

// CacheUser caches user data
func (r *Redis) CacheUser(ctx context.Context, userID uint, user interface{}, expiration time.Duration) error {
	key := fmt.Sprintf("%s%d", UserCacheKeyPrefix, userID)
	return r.Set(ctx, key, user, expiration)
}

// GetCachedUser retrieves cached user data
func (r *Redis) GetCachedUser(ctx context.Context, userID uint, dest interface{}) error {
	key := fmt.Sprintf("%s%d", UserCacheKeyPrefix, userID)
	return r.Get(ctx, key, dest)
}

// InvalidateUserCache removes user from cache
func (r *Redis) InvalidateUserCache(ctx context.Context, userID uint) error {
	key := fmt.Sprintf("%s%d", UserCacheKeyPrefix, userID)
	return r.Delete(ctx, key)
}

// SetRateLimit sets rate limit counter
func (r *Redis) SetRateLimit(ctx context.Context, identifier string, expiration time.Duration) (int64, error) {
	key := fmt.Sprintf("%s%s", RateLimitKeyPrefix, identifier)
	return r.IncrementWithExpiry(ctx, key, expiration)
}

// BlacklistToken adds a token to blacklist
func (r *Redis) BlacklistToken(ctx context.Context, token string, expiration time.Duration) error {
	key := fmt.Sprintf("%s%s", TokenBlacklistPrefix, token)
	return r.Set(ctx, key, true, expiration)
}

// IsTokenBlacklisted checks if a token is blacklisted
func (r *Redis) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	key := fmt.Sprintf("%s%s", TokenBlacklistPrefix, token)
	return r.Exists(ctx, key)
}
