package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/zhwjimmy/user-center/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New creates a new logger based on configuration
func New(cfg config.LoggingConfig) (*zap.Logger, error) {
	// Create logs directory if it doesn't exist
	if cfg.OutputPath != "" {
		dir := filepath.Dir(cfg.OutputPath)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}
	}

	// Configure log level
	var level zapcore.Level
	switch cfg.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// Configure encoder
	var encoderConfig zapcore.EncoderConfig
	var encoder zapcore.Encoder

	if cfg.Format == "json" {
		encoderConfig = zap.NewProductionEncoderConfig()
		encoderConfig.TimeKey = "time"
		encoderConfig.LevelKey = "level"
		encoderConfig.MessageKey = "message"
		encoderConfig.CallerKey = "caller"
		encoderConfig.StacktraceKey = "stacktrace"
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
		encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// Configure output
	var writeSyncer zapcore.WriteSyncer
	if cfg.OutputPath != "" {
		file, err := os.OpenFile(cfg.OutputPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
		writeSyncer = zapcore.AddSync(file)
	} else {
		writeSyncer = zapcore.AddSync(os.Stdout)
	}

	// Create core
	core := zapcore.NewCore(encoder, writeSyncer, level)

	// Create logger with options
	logger := zap.New(core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.Development(),
	)

	return logger, nil
}

// WithRequestID adds request ID to logger
func WithRequestID(logger *zap.Logger, requestID string) *zap.Logger {
	return logger.With(zap.String("request_id", requestID))
}

// WithUserID adds user ID to logger
func WithUserID(logger *zap.Logger, userID uint) *zap.Logger {
	return logger.With(zap.Uint("user_id", userID))
}

// WithFields adds multiple fields to logger
func WithFields(logger *zap.Logger, fields map[string]interface{}) *zap.Logger {
	zapFields := make([]zap.Field, 0, len(fields))
	for key, value := range fields {
		zapFields = append(zapFields, zap.Any(key, value))
	}
	return logger.With(zapFields...)
}
