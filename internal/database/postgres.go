package database

import (
	"fmt"
	"time"

	"github.com/your-org/user-center/internal/config"
	"github.com/your-org/user-center/internal/model"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// PostgreSQL represents PostgreSQL database connection
type PostgreSQL struct {
	DB *gorm.DB
}

// NewPostgreSQL creates a new PostgreSQL connection
func NewPostgreSQL(cfg *config.Config, zapLogger *zap.Logger) (*PostgreSQL, error) {
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

	return &PostgreSQL{DB: db}, nil
}

// Close closes the database connection
func (p *PostgreSQL) Close() error {
	sqlDB, err := p.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Health checks the database health
func (p *PostgreSQL) Health() error {
	sqlDB, err := p.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// GormZapWriter implements GORM logger interface for Zap
type GormZapWriter struct {
	logger *zap.Logger
}

// Printf implements the logger.Writer interface
func (w *GormZapWriter) Printf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	w.logger.Info(message)
}
