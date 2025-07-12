package service

import (
	"context"
	"fmt"

	"github.com/zhwjimmy/user-center/internal/dto"
	"github.com/zhwjimmy/user-center/internal/model"
	"github.com/zhwjimmy/user-center/internal/repository"
	"go.uber.org/zap"
)

// UserService handles user business logic
type UserService struct {
	userRepo repository.UserRepository
	logger   *zap.Logger
}

// NewUserService creates a new user service
func NewUserService(
	userRepo repository.UserRepository,
	logger *zap.Logger,
) *UserService {
	return &UserService{
		userRepo: userRepo,
		logger:   logger,
	}
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get user by ID",
			zap.String("user_id", id),
			zap.Error(err),
		)
		return nil, err
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		s.logger.Error("Failed to get user by email",
			zap.String("email", email),
			zap.Error(err),
		)
		return nil, err
	}

	return user, nil
}

// GetUserByUsername retrieves a user by username
func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		s.logger.Error("Failed to get user by username",
			zap.String("username", username),
			zap.Error(err),
		)
		return nil, err
	}

	return user, nil
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	// Check if user with email already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", user.Email)
	}

	// Check if user with username already exists
	existingUser, err = s.userRepo.GetByUsername(ctx, user.Username)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with username %s already exists", user.Username)
	}

	createdUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		s.logger.Error("Failed to create user",
			zap.String("email", user.Email),
			zap.String("username", user.Username),
			zap.Error(err),
		)
		return nil, err
	}

	s.logger.Info("User created successfully",
		zap.String("user_id", createdUser.ID),
		zap.String("email", createdUser.Email),
		zap.String("username", createdUser.Username),
	)

	return createdUser, nil
}

// UpdateUser updates user information
func (s *UserService) UpdateUser(ctx context.Context, id string, req *dto.UpdateUserRequest) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.FirstName != nil {
		user.FirstName = req.FirstName
	}
	if req.LastName != nil {
		user.LastName = req.LastName
	}
	if req.Avatar != nil {
		user.AvatarURL = req.Avatar
	}
	if req.Phone != nil {
		user.Phone = req.Phone
	}

	updatedUser, err := s.userRepo.Update(ctx, user)
	if err != nil {
		s.logger.Error("Failed to update user",
			zap.String("user_id", id),
			zap.Error(err),
		)
		return nil, err
	}

	s.logger.Info("User updated successfully",
		zap.String("user_id", updatedUser.ID),
	)

	return updatedUser, nil
}

// DeleteUser soft deletes a user
func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	err := s.userRepo.Delete(ctx, id)
	if err != nil {
		s.logger.Error("Failed to delete user",
			zap.String("user_id", id),
			zap.Error(err),
		)
		return err
	}

	s.logger.Info("User deleted successfully",
		zap.String("user_id", id),
	)

	return nil
}

// ListUsers retrieves users with pagination and filters
func (s *UserService) ListUsers(ctx context.Context, req *dto.UserListRequest) ([]*model.User, int64, error) {
	users, total, err := s.userRepo.List(ctx, req)
	if err != nil {
		s.logger.Error("Failed to list users",
			zap.Error(err),
		)
		return nil, 0, err
	}

	s.logger.Debug("Users listed successfully",
		zap.Int("count", len(users)),
		zap.Int64("total", total),
		zap.Int("page", req.Page),
		zap.Int("size", req.Size),
	)

	return users, total, nil
}

// UpdateUserStatus updates user status
func (s *UserService) UpdateUserStatus(ctx context.Context, id string, status model.UserStatus) (*model.User, error) {
	if !status.IsValid() {
		return nil, fmt.Errorf("invalid user status: %s", status)
	}

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update user status based on the status enum
	switch status {
	case model.UserStatusActive:
		user.IsActive = true
	case model.UserStatusInactive, model.UserStatusSuspended, model.UserStatusDeleted:
		user.IsActive = false
	}

	updatedUser, err := s.userRepo.Update(ctx, user)
	if err != nil {
		s.logger.Error("Failed to update user status",
			zap.String("user_id", id),
			zap.String("status", string(status)),
			zap.Error(err),
		)
		return nil, err
	}

	s.logger.Info("User status updated successfully",
		zap.String("user_id", updatedUser.ID),
		zap.String("status", string(status)),
	)

	return updatedUser, nil
}

// ActivateUser activates a user account
func (s *UserService) ActivateUser(ctx context.Context, id string) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.IsActive = true

	updatedUser, err := s.userRepo.Update(ctx, user)
	if err != nil {
		s.logger.Error("Failed to activate user",
			zap.String("user_id", id),
			zap.Error(err),
		)
		return nil, err
	}

	s.logger.Info("User activated successfully",
		zap.String("user_id", updatedUser.ID),
	)

	return updatedUser, nil
}

// DeactivateUser deactivates a user account
func (s *UserService) DeactivateUser(ctx context.Context, id string) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.IsActive = false

	updatedUser, err := s.userRepo.Update(ctx, user)
	if err != nil {
		s.logger.Error("Failed to deactivate user",
			zap.String("user_id", id),
			zap.Error(err),
		)
		return nil, err
	}

	s.logger.Info("User deactivated successfully",
		zap.String("user_id", updatedUser.ID),
	)

	return updatedUser, nil
}

// SearchUsers searches users by term
func (s *UserService) SearchUsers(ctx context.Context, term string, limit int) ([]*model.User, error) {
	users, err := s.userRepo.Search(ctx, term, limit)
	if err != nil {
		s.logger.Error("Failed to search users",
			zap.String("term", term),
			zap.Error(err),
		)
		return nil, err
	}

	s.logger.Debug("Users searched successfully",
		zap.String("term", term),
		zap.Int("count", len(users)),
	)

	return users, nil
}
