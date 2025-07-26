package asynq

import (
	"github.com/hibiken/asynq"
	"github.com/zhwjimmy/user-center/internal/config"
)

// AsynqConfig Asynq 配置
type AsynqConfig struct {
	Redis    RedisConfig    `mapstructure:"redis"`
	Queues   map[string]int `mapstructure:"queues"`
	Workers  int            `mapstructure:"workers"`
	LogLevel string         `mapstructure:"log_level"`
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// NewAsynqConfig 从应用配置创建 Asynq 配置
func NewAsynqConfig(cfg *config.Config) *AsynqConfig {
	return &AsynqConfig{
		Redis: RedisConfig{
			Addr:     cfg.Task.Redis.Addr,
			Password: cfg.Task.Redis.Password,
			DB:       cfg.Task.Redis.DB,
		},
		Queues: map[string]int{
			"critical": 6,
			"default":  3,
			"low":      1,
		},
		Workers:  cfg.Task.Workers,
		LogLevel: cfg.Task.LogLevel,
	}
}

// GetRedisClientOpt 获取 Redis 客户端选项
func (ac *AsynqConfig) GetRedisClientOpt() asynq.RedisClientOpt {
	return asynq.RedisClientOpt{
		Addr:     ac.Redis.Addr,
		Password: ac.Redis.Password,
		DB:       ac.Redis.DB,
	}
}

// GetServerConfig 获取服务器配置
func (ac *AsynqConfig) GetServerConfig() asynq.Config {
	return asynq.Config{
		Concurrency: ac.Workers,
		Queues:      ac.Queues,
		Logger:      nil, // 使用默认日志器
		LogLevel:    getLogLevel(ac.LogLevel),
	}
}

// getLogLevel 转换日志级别
func getLogLevel(level string) asynq.LogLevel {
	switch level {
	case "debug":
		return asynq.DebugLevel
	case "info":
		return asynq.InfoLevel
	case "warn":
		return asynq.WarnLevel
	case "error":
		return asynq.ErrorLevel
	default:
		return asynq.InfoLevel
	}
}
