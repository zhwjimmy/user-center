package infrastructure

import (
	"context"
	"fmt"
	"time"
)

// HealthStatus 健康状态
type HealthStatus struct {
	Overall   string            `json:"overall"`
	Timestamp time.Time         `json:"timestamp"`
	Services  map[string]Status `json:"services"`
}

// Status 服务状态
type Status struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// Health 检查所有基础设施服务的健康状态
func (m *Manager) Health(ctx context.Context) HealthStatus {
	status := HealthStatus{
		Timestamp: time.Now(),
		Services:  make(map[string]Status),
	}

	// 检查 PostgreSQL
	if err := m.postgres.Health(); err != nil {
		status.Services["postgresql"] = Status{
			Status:  "unhealthy",
			Message: err.Error(),
		}
	} else {
		status.Services["postgresql"] = Status{Status: "healthy"}
	}

	// 检查 MongoDB
	if err := m.mongodb.Health(ctx); err != nil {
		status.Services["mongodb"] = Status{
			Status:  "unhealthy",
			Message: err.Error(),
		}
	} else {
		status.Services["mongodb"] = Status{Status: "healthy"}
	}

	// 检查 Redis
	if err := m.redis.Health(ctx); err != nil {
		status.Services["redis"] = Status{
			Status:  "unhealthy",
			Message: err.Error(),
		}
	} else {
		status.Services["redis"] = Status{Status: "healthy"}
	}

	// 检查 Kafka
	if m.kafka != nil {
		status.Services["kafka"] = Status{Status: "healthy"}
	} else {
		status.Services["kafka"] = Status{
			Status:  "unhealthy",
			Message: "kafka service not initialized",
		}
	}

	// 确定整体状态
	overall := "healthy"
	for _, serviceStatus := range status.Services {
		if serviceStatus.Status == "unhealthy" {
			overall = "unhealthy"
			break
		}
	}
	status.Overall = overall

	return status
}

// Ready 检查服务是否就绪
func (m *Manager) Ready(ctx context.Context) error {
	status := m.Health(ctx)
	if status.Overall != "healthy" {
		return fmt.Errorf("infrastructure not ready: %+v", status.Services)
	}
	return nil
}
