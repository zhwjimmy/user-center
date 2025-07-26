package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zhwjimmy/user-center/internal/config"
	"go.uber.org/zap"
)

// redisImpl 实现 Cache 接口
type redisImpl struct {
	client *redis.Client
	logger *zap.Logger
}

// 确保 redisImpl 实现了 Cache 接口
var _ Cache = (*redisImpl)(nil)

// NewRedis 创建新的 Redis 连接
func NewRedis(cfg *config.Config, logger *zap.Logger) (Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.Addr,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.Info("Redis connected successfully",
		zap.String("addr", cfg.Redis.Addr),
		zap.Int("db", cfg.Redis.DB),
	)

	return &redisImpl{
		client: client,
		logger: logger,
	}, nil
}

// Close 关闭 Redis 连接
func (r *redisImpl) Close() error {
	return r.client.Close()
}

// Health 检查 Redis 健康状态
func (r *redisImpl) Health(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

// Set 存储带过期时间的值
func (r *redisImpl) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	if err := r.client.Set(ctx, key, data, expiration).Err(); err != nil {
		r.logger.Error("Failed to set cache",
			zap.String("key", key),
			zap.Error(err),
		)
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

// Get 从缓存中获取值
func (r *redisImpl) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := r.client.Get(ctx, key).Result()
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

// Delete 从缓存中删除键
func (r *redisImpl) Delete(ctx context.Context, key string) error {
	if err := r.client.Del(ctx, key).Err(); err != nil {
		r.logger.Error("Failed to delete cache",
			zap.String("key", key),
			zap.Error(err),
		)
		return fmt.Errorf("failed to delete cache: %w", err)
	}

	return nil
}

// Exists 检查键是否存在
func (r *redisImpl) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		r.logger.Error("Failed to check cache existence",
			zap.String("key", key),
			zap.Error(err),
		)
		return false, fmt.Errorf("failed to check cache existence: %w", err)
	}

	return result > 0, nil
}

// SetNX 仅在键不存在时设置值
func (r *redisImpl) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("failed to marshal value: %w", err)
	}

	result, err := r.client.SetNX(ctx, key, data, expiration).Result()
	if err != nil {
		r.logger.Error("Failed to set cache NX",
			zap.String("key", key),
			zap.Error(err),
		)
		return false, fmt.Errorf("failed to set cache NX: %w", err)
	}

	return result, nil
}

// Increment 递增计数器
func (r *redisImpl) Increment(ctx context.Context, key string) (int64, error) {
	result, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		r.logger.Error("Failed to increment counter",
			zap.String("key", key),
			zap.Error(err),
		)
		return 0, fmt.Errorf("failed to increment counter: %w", err)
	}

	return result, nil
}

// IncrementWithExpiry 递增带过期时间的计数器
func (r *redisImpl) IncrementWithExpiry(ctx context.Context, key string, expiration time.Duration) (int64, error) {
	pipe := r.client.Pipeline()
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

// SetExpiry 设置键的过期时间
func (r *redisImpl) SetExpiry(ctx context.Context, key string, expiration time.Duration) error {
	if err := r.client.Expire(ctx, key, expiration).Err(); err != nil {
		r.logger.Error("Failed to set expiry",
			zap.String("key", key),
			zap.Error(err),
		)
		return fmt.Errorf("failed to set expiry: %w", err)
	}

	return nil
}

// GetTTL 获取键的剩余生存时间
func (r *redisImpl) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		r.logger.Error("Failed to get TTL",
			zap.String("key", key),
			zap.Error(err),
		)
		return 0, fmt.Errorf("failed to get TTL: %w", err)
	}

	return ttl, nil
}

// Keys 返回匹配模式的所有键
func (r *redisImpl) Keys(ctx context.Context, pattern string) ([]string, error) {
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		r.logger.Error("Failed to get keys",
			zap.String("pattern", pattern),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get keys: %w", err)
	}

	return keys, nil
}
