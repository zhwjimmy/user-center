package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/your-org/user-center/internal/dto"
	"github.com/your-org/user-center/internal/model"
	"github.com/your-org/user-center/internal/service"
	"github.com/your-org/user-center/pkg/jwt"
	"go.uber.org/zap"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService *service.UserService
	authService *service.AuthService
	logger      *zap.Logger
}

// NewUserHandler creates a new user handler
func NewUserHandler(
	userService *service.UserService,
	authService *service.AuthService,
	logger *zap.Logger,
) *UserHandler {
	return &UserHandler{
		userService: userService,
		authService: authService,
		logger:      logger,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user with username, email, and password
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration request"
// @Success 201 {object} dto.RegisterResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid registration request", zap.Error(err))
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: err.Error(),
		})
		return
	}

	user, token, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Registration failed", zap.Error(err))

		// Check for specific errors
		if err.Error() == "user already exists" {
			c.JSON(http.StatusConflict, dto.ErrorResponse{
				Error:   "Conflict",
				Message: "User with this email or username already exists",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to register user",
		})
		return
	}

	c.JSON(http.StatusCreated, dto.RegisterResponse{
		User:    user.ToPublicUser(),
		Token:   token,
		Message: "User registered successfully",
	})
}

// Login handles user login
// @Summary User login
// @Description Authenticate user with email and password
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login request"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid login request", zap.Error(err))
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: err.Error(),
		})
		return
	}

	user, token, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Login failed", zap.Error(err))

		if err.Error() == "invalid credentials" {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error:   "Unauthorized",
				Message: "Invalid email or password",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to login",
		})
		return
	}

	c.JSON(http.StatusOK, dto.LoginResponse{
		User:    user.ToPublicUser(),
		Token:   token,
		Message: "Login successful",
	})
}

// GetUser handles getting user by ID
// @Summary Get user by ID
// @Description Get user information by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error("Invalid user ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid user ID",
		})
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get user", zap.Error(err))

		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Not Found",
				Message: "User not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to get user",
		})
		return
	}

	c.JSON(http.StatusOK, dto.UserResponse{
		User:    user.ToPublicUser(),
		Message: "User retrieved successfully",
	})
}

// GetCurrentUser handles getting current user information
// @Summary Get current user
// @Description Get current authenticated user information
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} dto.UserResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /users/me [get]
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "Invalid token",
		})
		return
	}

	userClaims := claims.(*jwt.Claims)
	user, err := h.userService.GetUserByID(c.Request.Context(), userClaims.UserID)
	if err != nil {
		h.logger.Error("Failed to get current user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to get user",
		})
		return
	}

	c.JSON(http.StatusOK, dto.UserResponse{
		User:    user.ToPublicUser(),
		Message: "User retrieved successfully",
	})
}

// UpdateUser handles updating user information
// @Summary Update user
// @Description Update current user information
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.UpdateUserRequest true "Update request"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /users/me [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "Invalid token",
		})
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid update request", zap.Error(err))
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: err.Error(),
		})
		return
	}

	userClaims := claims.(*jwt.Claims)
	user, err := h.userService.UpdateUser(c.Request.Context(), userClaims.UserID, &req)
	if err != nil {
		h.logger.Error("Failed to update user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to update user",
		})
		return
	}

	c.JSON(http.StatusOK, dto.UserResponse{
		User:    user.ToPublicUser(),
		Message: "User updated successfully",
	})
}

// ListUsers handles getting user list with pagination and filters
// @Summary List users
// @Description Get paginated list of users with optional filters
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param size query int false "Page size" default(10)
// @Param sort query string false "Sort field" default(created_at)
// @Param order query string false "Sort order (asc/desc)" default(desc)
// @Param search query string false "Search term"
// @Param status query string false "User status"
// @Param is_active query bool false "User active status"
// @Success 200 {object} dto.UserListResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	var req dto.UserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.logger.Error("Invalid list request", zap.Error(err))
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: err.Error(),
		})
		return
	}

	// Set defaults if not provided
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Size < 1 {
		req.Size = 10
	}
	if req.Sort == "" {
		req.Sort = "created_at"
	}
	if req.Order == "" {
		req.Order = "desc"
	}

	users, total, err := h.userService.ListUsers(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to list users", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to list users",
		})
		return
	}

	// Convert to public users
	publicUsers := make([]*model.PublicUser, len(users))
	for i, user := range users {
		publicUsers[i] = user.ToPublicUser()
	}

	// Calculate pagination
	totalPages := int(total) / req.Size
	if int(total)%req.Size > 0 {
		totalPages++
	}

	pagination := &dto.PaginationResponse{
		Page:       req.Page,
		Size:       req.Size,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    req.Page < totalPages,
		HasPrev:    req.Page > 1,
	}

	c.JSON(http.StatusOK, dto.UserListResponse{
		Users:      publicUsers,
		Pagination: pagination,
		Message:    "Users retrieved successfully",
	})
}

// ChangePassword handles password change
// @Summary Change password
// @Description Change current user password
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.ChangePasswordRequest true "Change password request"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /users/me/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "Invalid token",
		})
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid change password request", zap.Error(err))
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: err.Error(),
		})
		return
	}

	userClaims := claims.(*jwt.Claims)
	err := h.authService.ChangePassword(c.Request.Context(), userClaims.UserID, &req)
	if err != nil {
		h.logger.Error("Failed to change password", zap.Error(err))

		if err.Error() == "invalid old password" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Bad Request",
				Message: "Invalid old password",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to change password",
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Password changed successfully",
	})
}
