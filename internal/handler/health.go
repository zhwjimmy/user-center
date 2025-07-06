package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zhwjimmy/user-center/internal/cache"
	"github.com/zhwjimmy/user-center/internal/database"
	"github.com/zhwjimmy/user-center/internal/dto"
	"go.uber.org/zap"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	logger   *zap.Logger
	postgres *database.PostgreSQL
	mongodb  *database.MongoDB
	redis    *cache.Redis
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(
	logger *zap.Logger,
	postgres *database.PostgreSQL,
	mongodb *database.MongoDB,
	redis *cache.Redis,
) *HealthHandler {
	return &HealthHandler{
		logger:   logger,
		postgres: postgres,
		mongodb:  mongodb,
		redis:    redis,
	}
}

// Health handles health check requests
// @Summary Health check
// @Description Check the health status of the service and its dependencies
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} dto.HealthResponse
// @Failure 503 {object} dto.HealthResponse
// @Router /health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	checks := make(map[string]string)
	overallStatus := "healthy"

	// Check PostgreSQL
	if err := h.checkPostgreSQL(); err != nil {
		checks["postgresql"] = "unhealthy: " + err.Error()
		overallStatus = "unhealthy"
		h.logger.Error("PostgreSQL health check failed", zap.Error(err))
	} else {
		checks["postgresql"] = "healthy"
	}

	// Check MongoDB
	if err := h.checkMongoDB(); err != nil {
		checks["mongodb"] = "unhealthy: " + err.Error()
		overallStatus = "unhealthy"
		h.logger.Error("MongoDB health check failed", zap.Error(err))
	} else {
		checks["mongodb"] = "healthy"
	}

	// Check Redis
	if err := h.checkRedis(); err != nil {
		checks["redis"] = "unhealthy: " + err.Error()
		overallStatus = "unhealthy"
		h.logger.Error("Redis health check failed", zap.Error(err))
	} else {
		checks["redis"] = "healthy"
	}

	response := dto.HealthResponse{
		Status:    overallStatus,
		Version:   "1.0.0",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Checks:    checks,
	}

	statusCode := http.StatusOK
	if overallStatus == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}

// Ready handles readiness probe requests
// @Summary Readiness check
// @Description Check if the service is ready to serve requests
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} dto.HealthResponse
// @Failure 503 {object} dto.HealthResponse
// @Router /ready [get]
func (h *HealthHandler) Ready(c *gin.Context) {
	// For readiness, we check if all critical dependencies are available
	checks := make(map[string]string)
	overallStatus := "ready"

	// Check PostgreSQL (critical for user operations)
	if err := h.checkPostgreSQL(); err != nil {
		checks["postgresql"] = "not ready: " + err.Error()
		overallStatus = "not ready"
	} else {
		checks["postgresql"] = "ready"
	}

	// Check Redis (critical for caching and sessions)
	if err := h.checkRedis(); err != nil {
		checks["redis"] = "not ready: " + err.Error()
		overallStatus = "not ready"
	} else {
		checks["redis"] = "ready"
	}

	// MongoDB is not critical for basic operations, so we don't fail readiness for it
	if err := h.checkMongoDB(); err != nil {
		checks["mongodb"] = "degraded: " + err.Error()
	} else {
		checks["mongodb"] = "ready"
	}

	response := dto.HealthResponse{
		Status:    overallStatus,
		Version:   "1.0.0",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Checks:    checks,
	}

	statusCode := http.StatusOK
	if overallStatus == "not ready" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}

// Live handles liveness probe requests
// @Summary Liveness check
// @Description Check if the service is alive
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} dto.HealthResponse
// @Router /live [get]
func (h *HealthHandler) Live(c *gin.Context) {
	// Liveness check is simple - just return OK if the service is running
	response := dto.HealthResponse{
		Status:    "alive",
		Version:   "1.0.0",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Checks: map[string]string{
			"service": "alive",
		},
	}

	c.JSON(http.StatusOK, response)
}

// checkPostgreSQL checks PostgreSQL connectivity
func (h *HealthHandler) checkPostgreSQL() error {
	if h.postgres == nil {
		return fmt.Errorf("postgres client not initialized")
	}

	db, err := h.postgres.DB.DB()
	if err != nil {
		return err
	}

	return db.Ping()
}

// checkMongoDB checks MongoDB connectivity
func (h *HealthHandler) checkMongoDB() error {
	if h.mongodb == nil {
		return fmt.Errorf("mongodb client not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return h.mongodb.Client.Ping(ctx, nil)
}

// checkRedis checks Redis connectivity
func (h *HealthHandler) checkRedis() error {
	if h.redis == nil {
		return fmt.Errorf("redis client not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return h.redis.Client.Ping(ctx).Err()
}
