//go:build wireinject
// +build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/zhwjimmy/user-center/internal/cache"
	"github.com/zhwjimmy/user-center/internal/config"
	"github.com/zhwjimmy/user-center/internal/database"
	"github.com/zhwjimmy/user-center/internal/handler"
	"github.com/zhwjimmy/user-center/internal/kafka"
	kafkaConfig "github.com/zhwjimmy/user-center/internal/kafka/config"
	"github.com/zhwjimmy/user-center/internal/middleware"
	"github.com/zhwjimmy/user-center/internal/repository"
	"github.com/zhwjimmy/user-center/internal/server"
	"github.com/zhwjimmy/user-center/internal/service"
	"github.com/zhwjimmy/user-center/pkg/jwt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
)

type (
	RecoveryMiddleware  gin.HandlerFunc
	LoggerMiddleware    gin.HandlerFunc
	RequestIDMiddleware gin.HandlerFunc
	CORSMiddleware      gin.HandlerFunc
)

// provideLogger creates a new logger instance
func provideLogger(cfg *config.Config) (*zap.Logger, error) {
	config := zap.NewProductionConfig()

	// Set log level
	level, err := zapcore.ParseLevel(cfg.Logging.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}
	config.Level = zap.NewAtomicLevelAt(level)

	// Set output format
	if cfg.Logging.Format == "console" {
		config.Encoding = "console"
	} else {
		config.Encoding = "json"
	}

	// Set output path if specified
	if cfg.Logging.OutputPath != "" {
		config.OutputPaths = []string{cfg.Logging.OutputPath}
	}

	return config.Build()
}

// provideJWT creates a new JWT manager
func provideJWT(cfg *config.Config) *jwt.JWT {
	return jwt.NewJWT(cfg.JWT.Secret, cfg.JWT.Issuer, cfg.JWT.Expiry)
}

// provideCORSMiddleware creates a new CORS middleware
func provideCORSMiddleware(cfg *config.Config) middleware.CORSMiddleware {
	return middleware.CORSMiddleware(middleware.NewCORSMiddleware(cfg))
}

// provideRequestIDMiddleware creates a new request ID middleware
func provideRequestIDMiddleware() middleware.RequestIDMiddleware {
	return middleware.RequestIDMiddleware(middleware.NewRequestIDMiddleware())
}

// provideLoggerMiddleware creates a new logger middleware
func provideLoggerMiddleware(logger *zap.Logger) middleware.LoggerMiddleware {
	return middleware.LoggerMiddleware(middleware.NewLoggerMiddleware(logger))
}

// provideRecoveryMiddleware creates a new recovery middleware
func provideRecoveryMiddleware(logger *zap.Logger) middleware.RecoveryMiddleware {
	return middleware.RecoveryMiddleware(middleware.NewRecoveryMiddleware(logger))
}

// provideGormDB extracts *gorm.DB from *database.PostgreSQL
func provideGormDB(pg *database.PostgreSQL) *gorm.DB {
	return pg.DB
}

// provideServer creates a new server instance
func provideServer(
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
) *server.Server {
	return server.New(
		cfg,
		logger,
		userHandler,
		healthHandler,
		authMiddleware,
		corsMiddleware,
		rateLimitMiddleware,
		requestIDMiddleware,
		loggerMiddleware,
		recoveryMiddleware,
		kafkaService,
	)
}

// InitializeApp creates a new application instance
func InitializeApp() (*server.Server, error) {
	wire.Build(
		// Configuration
		config.Load,
		kafkaConfig.NewKafkaClientConfig,

		// Logger
		provideLogger,

		// JWT Manager
		provideJWT,

		// Database connections
		database.NewPostgreSQL,
		provideGormDB,
		database.NewMongoDB,
		cache.NewRedis,

		// Kafka
		kafka.NewKafkaService,

		// Repositories
		repository.NewUserRepository,

		// Services
		service.NewUserService,
		service.NewEventService,
		service.NewAuthService,

		// Handlers
		handler.NewUserHandler,
		handler.NewHealthHandler,

		// Middlewares
		middleware.NewAuthMiddleware,
		provideCORSMiddleware,
		middleware.NewRateLimitMiddleware,
		provideRequestIDMiddleware,
		provideLoggerMiddleware,
		provideRecoveryMiddleware,

		// Server
		provideServer,
	)
	return &server.Server{}, nil
}
