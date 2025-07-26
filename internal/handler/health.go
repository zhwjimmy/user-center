package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zhwjimmy/user-center/internal/dto"
	"github.com/zhwjimmy/user-center/internal/infrastructure"
	"go.uber.org/zap"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	logger *zap.Logger
	infra  *infrastructure.Manager
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(
	logger *zap.Logger,
	infra *infrastructure.Manager,
) *HealthHandler {
	return &HealthHandler{
		logger: logger,
		infra:  infra,
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
	status := h.infra.Health(c.Request.Context())

	checks := make(map[string]string)
	for service, serviceStatus := range status.Services {
		checks[service] = serviceStatus.Status
		if serviceStatus.Message != "" {
			checks[service] += ": " + serviceStatus.Message
		}
	}

	response := dto.HealthResponse{
		Status:    status.Overall,
		Version:   "1.0.0",
		Timestamp: status.Timestamp.Format(time.RFC3339),
		Checks:    checks,
	}

	statusCode := http.StatusOK
	if status.Overall != "healthy" {
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
	// 使用基础设施管理器的就绪检查
	if err := h.infra.Ready(c.Request.Context()); err != nil {
		h.logger.Error("Service not ready", zap.Error(err))
		c.JSON(http.StatusServiceUnavailable, dto.HealthResponse{
			Status:    "not ready",
			Version:   "1.0.0",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Checks: map[string]string{
				"service": "not ready: " + err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, dto.HealthResponse{
		Status:    "ready",
		Version:   "1.0.0",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Checks: map[string]string{
			"service": "ready",
		},
	})
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
