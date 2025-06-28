package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	Database   DatabaseConfig   `mapstructure:"database"`
	Redis      RedisConfig      `mapstructure:"redis"`
	Kafka      KafkaConfig      `mapstructure:"kafka"`
	JWT        JWTConfig        `mapstructure:"jwt"`
	Logging    LoggingConfig    `mapstructure:"logging"`
	Monitoring MonitoringConfig `mapstructure:"monitoring"`
	I18n       I18nConfig       `mapstructure:"i18n"`
	RateLimit  RateLimitConfig  `mapstructure:"rate_limit"`
	CORS       CORSConfig       `mapstructure:"cors"`
	Task       TaskConfig       `mapstructure:"task"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Mode            string        `mapstructure:"mode"` // debug, release, test
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Postgres PostgreSQLConfig `mapstructure:"postgres"`
	MongoDB  MongoDBConfig    `mapstructure:"mongodb"`
}

// PostgreSQLConfig holds PostgreSQL configuration
type PostgreSQLConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	User         string        `mapstructure:"user"`
	Password     string        `mapstructure:"password"`
	DBName       string        `mapstructure:"dbname"`
	SSLMode      string        `mapstructure:"sslmode"`
	MaxOpenConns int           `mapstructure:"max_open_conns"`
	MaxIdleConns int           `mapstructure:"max_idle_conns"`
	MaxLifetime  time.Duration `mapstructure:"max_lifetime"`
}

// MongoDBConfig holds MongoDB configuration
type MongoDBConfig struct {
	URI      string `mapstructure:"uri"`
	Database string `mapstructure:"database"`
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Addr         string `mapstructure:"addr"`
	Password     string `mapstructure:"password"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
}

// KafkaConfig holds Kafka configuration
type KafkaConfig struct {
	Brokers []string          `mapstructure:"brokers"`
	Topics  map[string]string `mapstructure:"topics"`
	GroupID string            `mapstructure:"group_id"`
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret string        `mapstructure:"secret"`
	Expiry time.Duration `mapstructure:"expiry"`
	Issuer string        `mapstructure:"issuer"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"` // json, console
	OutputPath string `mapstructure:"output_path"`
}

// MonitoringConfig holds monitoring configuration
type MonitoringConfig struct {
	Prometheus PrometheusConfig `mapstructure:"prometheus"`
	Tracing    TracingConfig    `mapstructure:"tracing"`
}

// PrometheusConfig holds Prometheus configuration
type PrometheusConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Port    int    `mapstructure:"port"`
	Path    string `mapstructure:"path"`
}

// TracingConfig holds tracing configuration
type TracingConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Endpoint string `mapstructure:"endpoint"`
	Service  string `mapstructure:"service"`
}

// I18nConfig holds internationalization configuration
type I18nConfig struct {
	DefaultLanguage string   `mapstructure:"default_language"`
	Languages       []string `mapstructure:"languages"`
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Rate    int    `mapstructure:"rate"`
	Burst   int    `mapstructure:"burst"`
	Store   string `mapstructure:"store"` // memory, redis
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowOrigins     []string `mapstructure:"allow_origins"`
	AllowMethods     []string `mapstructure:"allow_methods"`
	AllowHeaders     []string `mapstructure:"allow_headers"`
	ExposeHeaders    []string `mapstructure:"expose_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
	MaxAge           int      `mapstructure:"max_age"`
}

// TaskConfig holds async task configuration
type TaskConfig struct {
	Redis    RedisConfig `mapstructure:"redis"`
	Queues   []string    `mapstructure:"queues"`
	Workers  int         `mapstructure:"workers"`
	LogLevel string      `mapstructure:"log_level"`
}

// Load loads configuration from file and environment variables
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// Set environment variable prefix
	viper.SetEnvPrefix("USERCENTER")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Set default values
	setDefaults()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	// Server defaults
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.shutdown_timeout", "30s")

	// Database defaults
	viper.SetDefault("database.postgres.host", "localhost")
	viper.SetDefault("database.postgres.port", 5432)
	viper.SetDefault("database.postgres.user", "postgres")
	viper.SetDefault("database.postgres.password", "")
	viper.SetDefault("database.postgres.dbname", "usercenter")
	viper.SetDefault("database.postgres.sslmode", "disable")
	viper.SetDefault("database.postgres.max_open_conns", 25)
	viper.SetDefault("database.postgres.max_idle_conns", 10)
	viper.SetDefault("database.postgres.max_lifetime", "5m")

	viper.SetDefault("database.mongodb.uri", "mongodb://localhost:27017")
	viper.SetDefault("database.mongodb.database", "usercenter_logs")

	// Redis defaults
	viper.SetDefault("redis.addr", "localhost:6379")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.pool_size", 10)
	viper.SetDefault("redis.min_idle_conns", 5)

	// Kafka defaults
	viper.SetDefault("kafka.brokers", []string{"localhost:9092"})
	viper.SetDefault("kafka.topics.user_events", "user.events")
	viper.SetDefault("kafka.group_id", "usercenter")

	// JWT defaults
	viper.SetDefault("jwt.secret", "your-secret-key")
	viper.SetDefault("jwt.expiry", "24h")
	viper.SetDefault("jwt.issuer", "usercenter")

	// Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
	viper.SetDefault("logging.output_path", "logs/usercenter.log")

	// Monitoring defaults
	viper.SetDefault("monitoring.prometheus.enabled", true)
	viper.SetDefault("monitoring.prometheus.port", 9090)
	viper.SetDefault("monitoring.prometheus.path", "/metrics")

	viper.SetDefault("monitoring.tracing.enabled", true)
	viper.SetDefault("monitoring.tracing.endpoint", "http://localhost:14268/api/traces")
	viper.SetDefault("monitoring.tracing.service", "usercenter")

	// I18n defaults
	viper.SetDefault("i18n.default_language", "zh-CN")
	viper.SetDefault("i18n.languages", []string{"zh-CN", "en-US"})

	// Rate limit defaults
	viper.SetDefault("rate_limit.enabled", true)
	viper.SetDefault("rate_limit.rate", 100)
	viper.SetDefault("rate_limit.burst", 200)
	viper.SetDefault("rate_limit.store", "redis")

	// CORS defaults
	viper.SetDefault("cors.allow_origins", []string{"*"})
	viper.SetDefault("cors.allow_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	viper.SetDefault("cors.allow_headers", []string{"*"})
	viper.SetDefault("cors.expose_headers", []string{"X-Request-ID"})
	viper.SetDefault("cors.allow_credentials", true)
	viper.SetDefault("cors.max_age", 86400)

	// Task defaults
	viper.SetDefault("task.queues", []string{"default", "email", "notification"})
	viper.SetDefault("task.workers", 10)
	viper.SetDefault("task.log_level", "info")
}

// GetDSN returns the PostgreSQL DSN
func (c *PostgreSQLConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}
