package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zhwjimmy/user-center/internal/config"
	"github.com/zhwjimmy/user-center/internal/dto"
	"github.com/zhwjimmy/user-center/internal/infrastructure/cache"
	"go.uber.org/zap"
)

// RateLimitMiddleware handles rate limiting
type RateLimitMiddleware struct {
	redis  cache.Cache
	config config.RateLimitConfig
	logger *zap.Logger
}

// NewRateLimitMiddleware creates a new rate limit middleware
func NewRateLimitMiddleware(redis cache.Cache, cfg *config.Config, logger *zap.Logger) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		redis:  redis,
		config: cfg.RateLimit,
		logger: logger,
	}
}

// RateLimit applies rate limiting based on client IP
func (m *RateLimitMiddleware) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !m.config.Enabled {
			c.Next()
			return
		}

		// Get client identifier (IP address)
		clientIP := c.ClientIP()

		// Create rate limit key
		key := fmt.Sprintf("rate_limit:%s", clientIP)

		// Check rate limit
		allowed, err := m.checkRateLimit(c.Request.Context(), key)
		if err != nil {
			m.logger.Error("Rate limit check failed",
				zap.String("client_ip", clientIP),
				zap.Error(err),
			)
			// Allow request if rate limit check fails
			c.Next()
			return
		}

		if !allowed {
			m.logger.Warn("Rate limit exceeded",
				zap.String("client_ip", clientIP),
			)
			c.JSON(http.StatusTooManyRequests, dto.ErrorResponse{
				Error:   "Too Many Requests",
				Message: "Rate limit exceeded. Please try again later.",
				Code:    "RATE_LIMIT_EXCEEDED",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitByUser applies rate limiting based on authenticated user
func (m *RateLimitMiddleware) RateLimitByUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !m.config.Enabled {
			c.Next()
			return
		}

		// Get user ID from context (set by auth middleware)
		userID, exists := c.Get("user_id")
		if !exists {
			// Fall back to IP-based rate limiting
			m.RateLimit()(c)
			return
		}

		// Create rate limit key
		key := fmt.Sprintf("rate_limit:user:%v", userID)

		// Check rate limit
		allowed, err := m.checkRateLimit(c.Request.Context(), key)
		if err != nil {
			m.logger.Error("User rate limit check failed",
				zap.Any("user_id", userID),
				zap.Error(err),
			)
			// Allow request if rate limit check fails
			c.Next()
			return
		}

		if !allowed {
			m.logger.Warn("User rate limit exceeded",
				zap.Any("user_id", userID),
			)
			c.JSON(http.StatusTooManyRequests, dto.ErrorResponse{
				Error:   "Too Many Requests",
				Message: "Rate limit exceeded. Please try again later.",
				Code:    "RATE_LIMIT_EXCEEDED",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitCustom applies rate limiting with custom parameters
func (m *RateLimitMiddleware) RateLimitCustom(rate int, window time.Duration, keyFunc func(*gin.Context) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !m.config.Enabled {
			c.Next()
			return
		}

		// Get custom key
		key := keyFunc(c)

		// Check rate limit with custom parameters
		allowed, err := m.checkCustomRateLimit(c.Request.Context(), key, rate, window)
		if err != nil {
			m.logger.Error("Custom rate limit check failed",
				zap.String("key", key),
				zap.Error(err),
			)
			// Allow request if rate limit check fails
			c.Next()
			return
		}

		if !allowed {
			m.logger.Warn("Custom rate limit exceeded",
				zap.String("key", key),
			)
			c.JSON(http.StatusTooManyRequests, dto.ErrorResponse{
				Error:   "Too Many Requests",
				Message: "Rate limit exceeded. Please try again later.",
				Code:    "RATE_LIMIT_EXCEEDED",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// checkRateLimit checks if the request is within rate limit
func (m *RateLimitMiddleware) checkRateLimit(ctx context.Context, key string) (bool, error) {
	window := time.Minute // Default window

	// Increment counter
	count, err := m.redis.IncrementWithExpiry(ctx, key, window)
	if err != nil {
		return false, err
	}

	// Check if within rate limit
	return count <= int64(m.config.Rate), nil
}

// checkCustomRateLimit checks rate limit with custom parameters
func (m *RateLimitMiddleware) checkCustomRateLimit(ctx context.Context, key string, rate int, window time.Duration) (bool, error) {
	// Increment counter
	count, err := m.redis.IncrementWithExpiry(ctx, key, window)
	if err != nil {
		return false, err
	}

	// Check if within rate limit
	return count <= int64(rate), nil
}

// LoginRateLimit applies rate limiting specifically for login attempts
func (m *RateLimitMiddleware) LoginRateLimit() gin.HandlerFunc {
	return m.RateLimitCustom(5, 15*time.Minute, func(c *gin.Context) string {
		return fmt.Sprintf("rate_limit:login:%s", c.ClientIP())
	})
}

// RegistrationRateLimit applies rate limiting specifically for registration attempts
func (m *RateLimitMiddleware) RegistrationRateLimit() gin.HandlerFunc {
	return m.RateLimitCustom(3, 1*time.Hour, func(c *gin.Context) string {
		return fmt.Sprintf("rate_limit:register:%s", c.ClientIP())
	})
}

// PasswordResetRateLimit applies rate limiting specifically for password reset attempts
func (m *RateLimitMiddleware) PasswordResetRateLimit() gin.HandlerFunc {
	return m.RateLimitCustom(3, 1*time.Hour, func(c *gin.Context) string {
		return fmt.Sprintf("rate_limit:password_reset:%s", c.ClientIP())
	})
}
