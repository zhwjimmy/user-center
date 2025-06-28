package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/your-org/user-center/internal/config"
)

// NewCORSMiddleware creates a new CORS middleware
func NewCORSMiddleware(cfg *config.Config) gin.HandlerFunc {
	corsConfig := cors.Config{
		AllowOrigins:     cfg.CORS.AllowOrigins,
		AllowMethods:     cfg.CORS.AllowMethods,
		AllowHeaders:     cfg.CORS.AllowHeaders,
		ExposeHeaders:    cfg.CORS.ExposeHeaders,
		AllowCredentials: cfg.CORS.AllowCredentials,
		MaxAge:           time.Duration(cfg.CORS.MaxAge) * time.Second,
	}

	return cors.New(corsConfig)
}
