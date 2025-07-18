package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/zhwjimmy/user-center/internal/config"
	"github.com/zhwjimmy/user-center/internal/handler"
	"github.com/zhwjimmy/user-center/internal/kafka"
	"github.com/zhwjimmy/user-center/internal/middleware"
	"go.uber.org/zap"
)

// Server represents the HTTP server
type Server struct {
	*gin.Engine
	config       *config.Config
	logger       *zap.Logger
	kafkaService kafka.Service
}

// New creates a new server instance
func New(
	cfg *config.Config,
	logger *zap.Logger,
	userHandler *handler.UserHandler,
	healthHandler *handler.HealthHandler,
	authMiddleware *middleware.AuthMiddleware,
	corsMiddleware middleware.CORSMiddleware,
	rateLimitMiddleware *middleware.RateLimitMiddleware,
	requestIDMiddleware middleware.RequestIDMiddleware,
	loggerMiddleware middleware.LoggerMiddleware,
	recoveryMiddleware middleware.RecoveryMiddleware,
	kafkaService kafka.Service,
) *Server {
	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Create Gin engine
	r := gin.New()

	// Global middleware
	r.Use(gin.HandlerFunc(recoveryMiddleware))
	r.Use(gin.HandlerFunc(requestIDMiddleware))
	r.Use(gin.HandlerFunc(loggerMiddleware))
	r.Use(gin.HandlerFunc(corsMiddleware))

	// Health check routes (no rate limiting or auth)
	r.GET("/health", healthHandler.Health)
	r.GET("/ready", healthHandler.Ready)
	r.GET("/live", healthHandler.Live)

	// Swagger documentation
	if cfg.Server.Mode != "release" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// API routes
	api := r.Group("/api")
	v1 := api.Group("/v1")

	// Public routes with rate limiting
	public := v1.Group("/")
	public.Use(rateLimitMiddleware.RateLimit())
	{
		// User registration and login
		users := public.Group("/users")
		{
			users.POST("/register",
				rateLimitMiddleware.RegistrationRateLimit(),
				userHandler.Register,
			)
			users.POST("/login",
				rateLimitMiddleware.LoginRateLimit(),
				userHandler.Login,
			)
		}
	}

	// Protected routes (require authentication)
	protected := v1.Group("/")
	protected.Use(authMiddleware.RequireAuth())
	protected.Use(authMiddleware.RequireActiveUser())
	protected.Use(rateLimitMiddleware.RateLimitByUser())
	{
		// User management
		users := protected.Group("/users")
		{
			users.GET("/:id", userHandler.GetUser)
			users.GET("/", userHandler.ListUsers)
			users.GET("/me", userHandler.GetCurrentUser)
			users.PUT("/me", userHandler.UpdateUser)
			users.PUT("/me/password", userHandler.ChangePassword)
		}
	}

	// Admin routes (require admin privileges)
	admin := v1.Group("/admin")
	admin.Use(authMiddleware.RequireAuth())
	admin.Use(authMiddleware.RequireActiveUser())
	admin.Use(authMiddleware.AdminOnly())
	admin.Use(rateLimitMiddleware.RateLimitByUser())
	{
		// Admin user management
		adminUsers := admin.Group("/users")
		{
			adminUsers.GET("/", userHandler.ListUsers)
			adminUsers.GET("/:id", userHandler.GetUser)
			// Additional admin-only endpoints can be added here
		}
	}

	// Metrics endpoint for Prometheus
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return &Server{
		Engine:       r,
		config:       cfg,
		logger:       logger,
		kafkaService: kafkaService,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)
	s.logger.Info("Starting HTTP server",
		zap.String("address", addr),
		zap.String("mode", s.config.Server.Mode),
	)
	return s.Run(addr)
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down HTTP server")

	// Create HTTP server instance for graceful shutdown
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port),
		Handler: s.Engine,
	}

	return srv.Shutdown(ctx)
}

// GetLogger returns the logger instance
func (s *Server) GetLogger() *zap.Logger {
	return s.logger
}

// GetShutdownTimeout returns the shutdown timeout from config
func (s *Server) GetShutdownTimeout() time.Duration {
	return s.config.Server.ShutdownTimeout
}
