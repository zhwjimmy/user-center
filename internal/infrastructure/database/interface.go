package database

import (
	"context"

	"gorm.io/gorm"
)

// PostgreSQL 数据库接口
type PostgreSQL interface {
	DB() *gorm.DB
	Close() error
	Health() error
}

// MongoDB 数据库接口
type MongoDB interface {
	Close(ctx context.Context) error
	Health(ctx context.Context) error
	Collection(name string) interface{}
}
