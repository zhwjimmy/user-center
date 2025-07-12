package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents the user entity
type User struct {
	ID            string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Username      string         `json:"username" gorm:"uniqueIndex;type:varchar(50);not null"`
	Email         string         `json:"email" gorm:"uniqueIndex;type:varchar(255);not null"`
	PasswordHash  string         `json:"-" gorm:"column:password_hash;type:varchar(255);not null"`
	FirstName     *string        `json:"first_name,omitempty" gorm:"column:first_name;type:varchar(100)"`
	LastName      *string        `json:"last_name,omitempty" gorm:"column:last_name;type:varchar(100)"`
	Phone         *string        `json:"phone,omitempty" gorm:"type:varchar(20)"`
	AvatarURL     *string        `json:"avatar_url,omitempty" gorm:"column:avatar_url;type:text"`
	IsActive      bool           `json:"is_active" gorm:"column:is_active;default:true"`
	IsAdmin       bool           `json:"is_admin" gorm:"column:is_admin;default:false"`
	EmailVerified bool           `json:"email_verified" gorm:"column:email_verified;default:false"`
	PhoneVerified bool           `json:"phone_verified" gorm:"column:phone_verified;default:false"`
	LastLoginAt   *time.Time     `json:"last_login_at,omitempty" gorm:"column:last_login_at;type:timestamp with time zone"`
	CreatedAt     time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
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
	if u.ID == "" {
		u.ID = uuid.New().String()
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
		ID:            u.ID,
		Username:      u.Username,
		Email:         u.Email,
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		Phone:         u.Phone,
		AvatarURL:     u.AvatarURL,
		IsActive:      u.IsActive,
		IsAdmin:       u.IsAdmin,
		EmailVerified: u.EmailVerified,
		PhoneVerified: u.PhoneVerified,
		LastLoginAt:   u.LastLoginAt,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}

// JWT interface methods to avoid circular dependency
func (u *User) GetID() string {
	return u.ID
}

func (u *User) GetUsername() string {
	return u.Username
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) GetStatus() string {
	if u.IsActive {
		return "active"
	}
	return "inactive"
}

// PublicUser represents public user information (without sensitive fields)
type PublicUser struct {
	ID            string     `json:"id"`
	Username      string     `json:"username"`
	Email         string     `json:"email"`
	FirstName     *string    `json:"first_name,omitempty"`
	LastName      *string    `json:"last_name,omitempty"`
	Phone         *string    `json:"phone,omitempty"`
	AvatarURL     *string    `json:"avatar_url,omitempty"`
	IsActive      bool       `json:"is_active"`
	IsAdmin       bool       `json:"is_admin"`
	EmailVerified bool       `json:"email_verified"`
	PhoneVerified bool       `json:"phone_verified"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
