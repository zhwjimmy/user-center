// Package database provides database interfaces and implementations for the user-center application.
// It defines common Database interfaces that can be implemented by different database backends
// such as PostgreSQL, MongoDB, etc.
package database

import (
	"context"

	"gorm.io/gorm"
)

// Database 通用数据库接口
type Database interface {
	GetDB() interface{} // 返回通用的数据库连接
	Close() error
	Health() error
}

// PostgresDB PostgreSQL 特定接口
type PostgresDB interface {
	Database
	DB() *gorm.DB // PostgreSQL 特定的 GORM 方法
}

// MongoDB 数据库接口
type MongoDB interface {
	Close(ctx context.Context) error
	Health(ctx context.Context) error
	Collection(name string) interface{}
}
