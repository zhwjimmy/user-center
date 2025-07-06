package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/zhwjimmy/user-center/internal/dto"
	"github.com/zhwjimmy/user-center/internal/model"
	"gorm.io/gorm"
)

// UserRepository defines user data access interface
//go:generate mockgen -destination=../mock/user_repository_mock.go -package=mock github.com/zhwjimmy/user-center/internal/repository UserRepository
// 注意：上面go:generate用于mockgen自动生成

type UserRepository interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)
	GetByID(ctx context.Context, id uint) (*model.User, error)
	GetByUUID(ctx context.Context, uuid string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	Update(ctx context.Context, user *model.User) (*model.User, error)
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, req *dto.UserListRequest) ([]*model.User, int64, error)
	Search(ctx context.Context, term string, limit int) ([]*model.User, error)
	GetByIDs(ctx context.Context, ids []uint) ([]*model.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	UpdateStatus(ctx context.Context, id uint, status model.UserStatus) error
	UpdateActiveStatus(ctx context.Context, id uint, isActive bool) error
	GetActiveUsers(ctx context.Context) ([]*model.User, error)
	GetUsersByStatus(ctx context.Context, status model.UserStatus) ([]*model.User, error)
	CountUsers(ctx context.Context) (int64, error)
	CountActiveUsers(ctx context.Context) (int64, error)
}

// userRepository is the concrete implementation
// of UserRepository interface
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return &user, nil
}

// GetByUUID retrieves a user by UUID
func (r *userRepository) GetByUUID(ctx context.Context, uuid string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("uuid = ?", uuid).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by UUID: %w", err)
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &user, nil
}

// GetByUsername retrieves a user by username
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	return &user, nil
}

// Update updates a user
func (r *userRepository) Update(ctx context.Context, user *model.User) (*model.User, error) {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	return user, nil
}

// Delete soft deletes a user
func (r *userRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.User{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// List retrieves users with pagination and filters
func (r *userRepository) List(ctx context.Context, req *dto.UserListRequest) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	query := r.db.WithContext(ctx).Model(&model.User{})

	// Apply filters
	if req.Search != "" {
		searchTerm := "%" + strings.ToLower(req.Search) + "%"
		query = query.Where(
			"LOWER(username) LIKE ? OR LOWER(email) LIKE ? OR LOWER(first_name) LIKE ? OR LOWER(last_name) LIKE ?",
			searchTerm, searchTerm, searchTerm, searchTerm,
		)
	}

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	if req.IsActive != nil {
		query = query.Where("is_active = ?", *req.IsActive)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Apply sorting
	orderClause := fmt.Sprintf("%s %s", req.Sort, strings.ToUpper(req.Order))
	query = query.Order(orderClause)

	// Apply pagination
	offset := (req.Page - 1) * req.Size
	query = query.Offset(offset).Limit(req.Size)

	// Execute query
	if err := query.Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	return users, total, nil
}

// Search searches users by term
func (r *userRepository) Search(ctx context.Context, term string, limit int) ([]*model.User, error) {
	var users []*model.User
	searchTerm := "%" + strings.ToLower(term) + "%"

	query := r.db.WithContext(ctx).Where(
		"LOWER(username) LIKE ? OR LOWER(email) LIKE ? OR LOWER(first_name) LIKE ? OR LOWER(last_name) LIKE ?",
		searchTerm, searchTerm, searchTerm, searchTerm,
	).Limit(limit)

	if err := query.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}

	return users, nil
}

// GetByIDs retrieves multiple users by IDs
func (r *userRepository) GetByIDs(ctx context.Context, ids []uint) ([]*model.User, error) {
	var users []*model.User
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to get users by IDs: %w", err)
	}
	return users, nil
}

// ExistsByEmail checks if a user exists by email
func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check user existence by email: %w", err)
	}
	return count > 0, nil
}

// ExistsByUsername checks if a user exists by username
func (r *userRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check user existence by username: %w", err)
	}
	return count > 0, nil
}

// UpdateStatus updates user status
func (r *userRepository) UpdateStatus(ctx context.Context, id uint, status model.UserStatus) error {
	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Update("status", status).Error; err != nil {
		return fmt.Errorf("failed to update user status: %w", err)
	}
	return nil
}

// UpdateActiveStatus updates user active status
func (r *userRepository) UpdateActiveStatus(ctx context.Context, id uint, isActive bool) error {
	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Update("is_active", isActive).Error; err != nil {
		return fmt.Errorf("failed to update user active status: %w", err)
	}
	return nil
}

// GetActiveUsers retrieves all active users
func (r *userRepository) GetActiveUsers(ctx context.Context) ([]*model.User, error) {
	var users []*model.User
	if err := r.db.WithContext(ctx).Where("is_active = ? AND status = ?", true, model.UserStatusActive).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to get active users: %w", err)
	}
	return users, nil
}

// GetUsersByStatus retrieves users by status
func (r *userRepository) GetUsersByStatus(ctx context.Context, status model.UserStatus) ([]*model.User, error) {
	var users []*model.User
	if err := r.db.WithContext(ctx).Where("status = ?", status).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to get users by status: %w", err)
	}
	return users, nil
}

// CountUsers returns the total number of users
func (r *userRepository) CountUsers(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return count, nil
}

// CountActiveUsers returns the number of active users
func (r *userRepository) CountActiveUsers(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("is_active = ? AND status = ?", true, model.UserStatusActive).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count active users: %w", err)
	}
	return count, nil
}
