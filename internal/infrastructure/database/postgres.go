package database

import (
	"fmt"
	"time"

	"github.com/zhwjimmy/user-center/internal/config"
	"github.com/zhwjimmy/user-center/internal/model"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// postgreSQL 实现 PostgreSQL 接口
type postgreSQL struct {
	db *gorm.DB
}

// 确保 postgreSQL 实现了 PostgreSQL 接口
var _ PostgreSQL = (*postgreSQL)(nil)

// NewPostgreSQL 创建新的 PostgreSQL 连接
func NewPostgreSQL(cfg *config.Config, zapLogger *zap.Logger) (PostgreSQL, error) {
	dsn := cfg.Database.Postgres.GetDSN()

	// Configure GORM logger
	gormLogger := logger.New(
		&GormZapWriter{logger: zapLogger},
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	// Open database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	// Get underlying sql.DB for connection pool configuration
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(cfg.Database.Postgres.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.Postgres.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.Postgres.MaxLifetime)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}

	// Auto migrate models
	if err := db.AutoMigrate(&model.User{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	zapLogger.Info("PostgreSQL connected successfully",
		zap.String("host", cfg.Database.Postgres.Host),
		zap.Int("port", cfg.Database.Postgres.Port),
		zap.String("database", cfg.Database.Postgres.DBName),
	)

	return &postgreSQL{db: db}, nil
}

// DB 返回 GORM DB 实例
func (p *postgreSQL) DB() *gorm.DB {
	return p.db
}

// Close 关闭数据库连接
func (p *postgreSQL) Close() error {
	sqlDB, err := p.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Health 检查数据库健康状态
func (p *postgreSQL) Health() error {
	sqlDB, err := p.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// GormZapWriter 实现 GORM logger 接口
type GormZapWriter struct {
	logger *zap.Logger
}

// Printf 实现 logger.Writer 接口
func (w *GormZapWriter) Printf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	w.logger.Info(message)
}
