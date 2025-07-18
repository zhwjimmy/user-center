package service

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/zhwjimmy/user-center/internal/dto"
	"github.com/zhwjimmy/user-center/internal/model"
	"github.com/zhwjimmy/user-center/pkg/jwt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication business logic
type AuthService struct {
	userService  *UserService
	eventService *EventService // New
	jwtManager   *jwt.JWT
	logger       *zap.Logger
}

// NewAuthService creates a new auth service
func NewAuthService(
	userService *UserService,
	eventService *EventService, // New
	jwtManager *jwt.JWT,
	logger *zap.Logger,
) *AuthService {
	return &AuthService{
		userService:  userService,
		eventService: eventService, // New
		jwtManager:   jwtManager,
		logger:       logger,
	}
}

// Register handles user registration
func (s *AuthService) Register(ctx context.Context, req *dto.RegisterRequest) (*model.User, string, error) {
	// Hash password
	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		s.logger.Error("Failed to hash password", zap.Error(err))
		return nil, "", fmt.Errorf("failed to process password")
	}

	// Create user model
	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Phone:        req.Phone,
		IsActive:     true,
	}

	// Create user
	createdUser, err := s.userService.CreateUser(ctx, user)
	if err != nil {
		s.logger.Error("Failed to create user during registration",
			zap.String("email", req.Email),
			zap.String("username", req.Username),
			zap.Error(err),
		)
		return nil, "", fmt.Errorf("user already exists")
	}

	// Generate JWT token
	token, err := s.jwtManager.GenerateToken(createdUser)
	if err != nil {
		s.logger.Error("Failed to generate token after registration",
			zap.String("user_id", createdUser.ID),
			zap.Error(err),
		)
		return nil, "", fmt.Errorf("failed to generate token")
	}

	// Publish user registration event
	if err := s.eventService.PublishUserRegisteredEvent(ctx, createdUser); err != nil {
		s.logger.Error("Failed to publish user registered event",
			zap.String("user_id", createdUser.ID),
			zap.Error(err),
		)
		// Do not return error to avoid affecting main business flow
	}

	s.logger.Info("User registered successfully",
		zap.String("user_id", createdUser.ID),
		zap.String("email", createdUser.Email),
		zap.String("username", createdUser.Username),
	)

	return createdUser, token, nil
}

// Login handles user login
func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequest) (*model.User, string, error) {
	// Get user by email
	user, err := s.userService.GetUserByEmail(ctx, req.Email)
	if err != nil {
		s.logger.Warn("Login attempt with non-existent email",
			zap.String("email", req.Email),
		)
		return nil, "", fmt.Errorf("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		s.logger.Warn("Login attempt with inactive user",
			zap.String("user_id", user.ID),
			zap.String("email", req.Email),
			zap.Bool("is_active", user.IsActive),
		)
		return nil, "", fmt.Errorf("account is inactive")
	}

	// Verify password
	if !s.verifyPassword(req.Password, user.PasswordHash) {
		s.logger.Warn("Login attempt with invalid password",
			zap.String("user_id", user.ID),
			zap.String("email", req.Email),
		)
		return nil, "", fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	token, err := s.jwtManager.GenerateToken(user)
	if err != nil {
		s.logger.Error("Failed to generate token after login",
			zap.String("user_id", user.ID),
			zap.Error(err),
		)
		return nil, "", fmt.Errorf("failed to generate token")
	}

	// Publish user login event
	ipAddress := s.getClientIP(ctx)
	userAgent := s.getUserAgent(ctx)
	if err := s.eventService.PublishUserLoggedInEvent(ctx, user, ipAddress, userAgent); err != nil {
		s.logger.Error("Failed to publish user logged in event",
			zap.String("user_id", user.ID),
			zap.Error(err),
		)
		// Do not return error to avoid affecting main business flow
	}

	s.logger.Info("User logged in successfully",
		zap.String("user_id", user.ID),
		zap.String("email", user.Email),
	)

	return user, token, nil
}

// ChangePassword handles password change
func (s *AuthService) ChangePassword(ctx context.Context, userID string, req *dto.ChangePasswordRequest) error {
	// Get user
	user, err := s.userService.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verify old password
	if !s.verifyPassword(req.OldPassword, user.PasswordHash) {
		s.logger.Warn("Invalid old password in change password request",
			zap.String("user_id", userID),
		)
		return fmt.Errorf("invalid old password")
	}

	// Hash new password
	hashedPassword, err := s.hashPassword(req.NewPassword)
	if err != nil {
		s.logger.Error("Failed to hash new password",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to process new password")
	}

	// Update password
	user.PasswordHash = hashedPassword
	_, err = s.userService.userRepo.Update(ctx, user)
	if err != nil {
		s.logger.Error("Failed to update password",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to update password")
	}

	// Publish user password changed event
	ipAddress := s.getClientIP(ctx)
	if err := s.eventService.PublishUserPasswordChangedEvent(ctx, user, ipAddress); err != nil {
		s.logger.Error("Failed to publish user password changed event",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		// Do not return error to avoid affecting main business flow
	}

	s.logger.Info("Password changed successfully",
		zap.String("user_id", userID),
	)

	return nil
}

// RefreshToken generates a new token from an existing token
func (s *AuthService) RefreshToken(ctx context.Context, tokenString string) (string, error) {
	// Validate existing token
	claims, err := s.jwtManager.ValidateToken(tokenString)
	if err != nil {
		s.logger.Warn("Invalid token in refresh request", zap.Error(err))
		return "", fmt.Errorf("invalid token")
	}

	// Get user to ensure they still exist and are active
	user, err := s.userService.GetUserByID(ctx, claims.UserID)
	if err != nil {
		s.logger.Warn("User not found during token refresh",
			zap.String("user_id", claims.UserID),
		)
		return "", fmt.Errorf("user not found")
	}

	// Check if user is still active
	if !user.IsActive {
		s.logger.Warn("Token refresh attempt for inactive user",
			zap.String("user_id", user.ID),
			zap.Bool("is_active", user.IsActive),
		)
		return "", fmt.Errorf("account is inactive")
	}

	// Generate new token
	newToken, err := s.jwtManager.GenerateToken(user)
	if err != nil {
		s.logger.Error("Failed to generate new token during refresh",
			zap.String("user_id", user.ID),
			zap.Error(err),
		)
		return "", fmt.Errorf("failed to generate token")
	}

	s.logger.Info("Token refreshed successfully",
		zap.String("user_id", user.ID),
	)

	return newToken, nil
}

// ValidateToken validates a JWT token and returns user claims
func (s *AuthService) ValidateToken(tokenString string) (*jwt.Claims, error) {
	return s.jwtManager.ValidateToken(tokenString)
}

// hashPassword hashes a password using bcrypt
func (s *AuthService) hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// verifyPassword verifies a password against its hash
func (s *AuthService) verifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ForgotPassword handles password reset request (placeholder for future implementation)
func (s *AuthService) ForgotPassword(ctx context.Context, email string) error {
	// This is a placeholder for forgot password functionality
	// In a real implementation, this would:
	// 1. Generate a password reset token
	// 2. Store it in cache/database with expiry
	// 3. Send password reset email

	user, err := s.userService.GetUserByEmail(ctx, email)
	if err != nil {
		// Don't reveal if email exists or not
		s.logger.Info("Password reset requested", zap.String("email", email))
		return nil
	}

	s.logger.Info("Password reset requested for existing user",
		zap.String("user_id", user.ID),
		zap.String("email", email),
	)

	// TODO: Implement password reset logic
	return nil
}

// ResetPassword handles password reset with token (placeholder for future implementation)
func (s *AuthService) ResetPassword(ctx context.Context, token, newPassword string) error {
	// This is a placeholder for reset password functionality
	// In a real implementation, this would:
	// 1. Validate the reset token
	// 2. Get user ID from token
	// 3. Update user password
	// 4. Invalidate the reset token

	s.logger.Info("Password reset attempted", zap.String("token", token))

	// TODO: Implement password reset logic
	return fmt.Errorf("password reset not implemented")
}

// 辅助方法
func (s *AuthService) getClientIP(ctx context.Context) string {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		return ginCtx.ClientIP()
	}
	return ""
}

func (s *AuthService) getUserAgent(ctx context.Context) string {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		return ginCtx.GetHeader("User-Agent")
	}
	return ""
}
