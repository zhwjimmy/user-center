package middleware

import (
	"time"

	"github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// NewLoggerMiddleware creates a new logger middleware
func NewLoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return ginzap.GinzapWithConfig(logger, &ginzap.Config{
		TimeFormat: time.RFC3339,
		UTC:        true,
		SkipPaths:  []string{"/health", "/ready", "/live"},
	})
}

// NewRecoveryMiddleware creates a new recovery middleware
func NewRecoveryMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return ginzap.RecoveryWithZap(logger, true)
}
