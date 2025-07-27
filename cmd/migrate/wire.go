//go:build wireinject
// +build wireinject

// Package main is the entry point for the UserCenter migration tool
package main

import (
	"github.com/google/wire"
	"github.com/zhwjimmy/user-center/internal/config"
	"github.com/zhwjimmy/user-center/internal/infrastructure/database"
	"go.uber.org/zap"
)

// MigrationApp 迁移应用结构
type MigrationApp struct {
	config *config.Config
	logger *zap.Logger
	db     database.PostgresDB
}

// GetLogger 获取日志器
func (app *MigrationApp) GetLogger() *zap.Logger {
	return app.logger
}

// InitializeMigrationApp 初始化迁移应用
func InitializeMigrationApp(configPath string) (*MigrationApp, error) {
	wire.Build(
		provideConfig,
		provideLogger,
		database.NewPostgreSQL,
		wire.Struct(new(MigrationApp), "*"),
	)
	return &MigrationApp{}, nil
}

// provideConfig 提供配置
func provideConfig(configPath string) (*config.Config, error) {
	return config.Load()
}

// provideLogger 提供日志器
func provideLogger(cfg *config.Config) (*zap.Logger, error) {
	// 创建日志器配置
	logConfig := zap.NewProductionConfig()

	// 根据配置设置日志级别
	switch cfg.Logging.Level {
	case "debug":
		logConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		logConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		logConfig.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		logConfig.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		logConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	// 设置日志格式
	if cfg.Logging.Format == "console" {
		logConfig.Encoding = "console"
	} else {
		logConfig.Encoding = "json"
	}

	// 设置输出路径
	if cfg.Logging.OutputPath != "" {
		logConfig.OutputPaths = []string{cfg.Logging.OutputPath}
	}

	return logConfig.Build()
}
