package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents the user entity
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	UUID      string         `json:"uuid" gorm:"uniqueIndex;type:varchar(36);not null"`
	Username  string         `json:"username" gorm:"uniqueIndex;type:varchar(50);not null"`
	Email     string         `json:"email" gorm:"uniqueIndex;type:varchar(100);not null"`
	Password  string         `json:"-" gorm:"type:varchar(255);not null"`
	FirstName *string        `json:"first_name,omitempty" gorm:"type:varchar(50)"`
	LastName  *string        `json:"last_name,omitempty" gorm:"type:varchar(50)"`
	Avatar    *string        `json:"avatar,omitempty" gorm:"type:varchar(255)"`
	Phone     *string        `json:"phone,omitempty" gorm:"type:varchar(20)"`
	Status    UserStatus     `json:"status" gorm:"type:varchar(20);default:'active'"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// UserStatus represents user status
type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusDeleted   UserStatus = "deleted"
)

// BeforeCreate generates UUID before creating user
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.UUID == "" {
		u.UUID = uuid.New().String()
	}
	return nil
}

// TableName returns the table name for User model
func (User) TableName() string {
	return "users"
}

// IsValidStatus checks if the status is valid
func (s UserStatus) IsValid() bool {
	switch s {
	case UserStatusActive, UserStatusInactive, UserStatusSuspended, UserStatusDeleted:
		return true
	default:
		return false
	}
}

// ToPublicUser converts User to PublicUser (without sensitive fields)
func (u *User) ToPublicUser() *PublicUser {
	return &PublicUser{
		ID:        u.ID,
		UUID:      u.UUID,
		Username:  u.Username,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Avatar:    u.Avatar,
		Phone:     u.Phone,
		Status:    u.Status,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// JWT interface methods to avoid circular dependency
func (u *User) GetID() uint {
	return u.ID
}

func (u *User) GetUsername() string {
	return u.Username
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) GetStatus() string {
	return string(u.Status)
}

// PublicUser represents public user information (without sensitive fields)
type PublicUser struct {
	ID        uint       `json:"id"`
	UUID      string     `json:"uuid"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	FirstName *string    `json:"first_name,omitempty"`
	LastName  *string    `json:"last_name,omitempty"`
	Avatar    *string    `json:"avatar,omitempty"`
	Phone     *string    `json:"phone,omitempty"`
	Status    UserStatus `json:"status"`
	IsActive  bool       `json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
