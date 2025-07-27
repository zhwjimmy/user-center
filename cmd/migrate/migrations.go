package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
)

// MigrateUp 执行向上迁移
func (app *MigrationApp) MigrateUp(ctx context.Context) error {
	app.logger.Info("Starting database migration up")

	// 构建 goose 命令
	cmd := app.buildGooseCommand("up")

	// 执行迁移
	if err := app.executeGooseCommand(cmd); err != nil {
		app.logger.Error("Migration up failed", zap.Error(err))
		return fmt.Errorf("failed to execute migration up: %w", err)
	}

	app.logger.Info("Database migration up completed successfully")
	return nil
}

// MigrateDown 执行向下迁移
func (app *MigrationApp) MigrateDown(ctx context.Context, steps int) error {
	app.logger.Info("Starting database migration down", zap.Int("steps", steps))

	// 构建 goose 命令
	cmd := app.buildGooseCommand("down")
	if steps > 1 {
		cmd = append(cmd, "-steps", fmt.Sprintf("%d", steps))
	}

	// 执行迁移
	if err := app.executeGooseCommand(cmd); err != nil {
		app.logger.Error("Migration down failed", zap.Error(err))
		return fmt.Errorf("failed to execute migration down: %w", err)
	}

	app.logger.Info("Database migration down completed successfully")
	return nil
}

// MigrateStatus 查看迁移状态
func (app *MigrationApp) MigrateStatus(ctx context.Context) error {
	app.logger.Info("Checking migration status")

	// 构建 goose 命令
	cmd := app.buildGooseCommand("status")

	// 执行状态检查
	if err := app.executeGooseCommand(cmd); err != nil {
		app.logger.Error("Migration status check failed", zap.Error(err))
		return fmt.Errorf("failed to check migration status: %w", err)
	}

	return nil
}

// CreateMigration 创建新的迁移文件
func (app *MigrationApp) CreateMigration(ctx context.Context, name string) error {
	app.logger.Info("Creating new migration file", zap.String("name", name))

	// 构建 goose 命令
	cmd := app.buildGooseCommand("create", name, "sql")

	// 执行创建
	if err := app.executeGooseCommand(cmd); err != nil {
		app.logger.Error("Failed to create migration file", zap.Error(err))
		return fmt.Errorf("failed to create migration file: %w", err)
	}

	app.logger.Info("Migration file created successfully", zap.String("name", name))
	return nil
}

// buildGooseCommand 构建 goose 命令
func (app *MigrationApp) buildGooseCommand(args ...string) []string {
	// 获取迁移目录的绝对路径
	migrationsDir, err := filepath.Abs("migrations")
	if err != nil {
		app.logger.Error("Failed to get migrations directory path", zap.Error(err))
		migrationsDir = "migrations"
	}

	// 构建数据库连接字符串
	dsn := app.buildDSN()

	// 构建完整的命令
	cmd := []string{
		"goose",
		"-dir", migrationsDir,
		"postgres",
		dsn,
	}

	// 添加其他参数
	cmd = append(cmd, args...)

	return cmd
}

// buildDSN 构建数据库连接字符串
func (app *MigrationApp) buildDSN() string {
	pg := app.config.Database.Postgres

	// 构建 PostgreSQL 连接字符串
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		pg.Host,
		pg.Port,
		pg.User,
		pg.Password,
		pg.DBName,
		pg.SSLMode,
	)

	return dsn
}

// executeGooseCommand 执行 goose 命令
func (app *MigrationApp) executeGooseCommand(cmd []string) error {
	app.logger.Debug("Executing goose command", zap.Strings("command", cmd))

	// 创建命令
	execCmd := exec.CommandContext(context.Background(), cmd[0], cmd[1:]...)

	// 设置环境变量
	execCmd.Env = os.Environ()

	// 设置输出
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	// 执行命令
	if err := execCmd.Run(); err != nil {
		return fmt.Errorf("goose command failed: %w", err)
	}

	return nil
}

// validateGooseInstallation 验证 goose 是否已安装
func (app *MigrationApp) validateGooseInstallation() error {
	cmd := exec.Command("goose", "-version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("goose is not installed or not in PATH: %w", err)
	}
	return nil
}

// validateMigrationsDirectory 验证迁移目录是否存在
func (app *MigrationApp) validateMigrationsDirectory() error {
	if _, err := os.Stat("migrations"); os.IsNotExist(err) {
		return fmt.Errorf("migrations directory does not exist: %w", err)
	}
	return nil
}

// getMigrationFiles 获取迁移文件列表
func (app *MigrationApp) getMigrationFiles() ([]string, error) {
	files, err := os.ReadDir("migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var migrationFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}

	return migrationFiles, nil
}
