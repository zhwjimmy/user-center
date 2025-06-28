package dto

import "github.com/your-org/user-center/internal/model"

// RegisterRequest represents user registration request
type RegisterRequest struct {
	Username  string  `json:"username" binding:"required,min=3,max=50" example:"testuser"`
	Email     string  `json:"email" binding:"required,email,max=100" example:"test@example.com"`
	Password  string  `json:"password" binding:"required,min=8,max=50" example:"securepassword123"`
	FirstName *string `json:"first_name,omitempty" binding:"omitempty,max=50" example:"John"`
	LastName  *string `json:"last_name,omitempty" binding:"omitempty,max=50" example:"Doe"`
	Phone     *string `json:"phone,omitempty" binding:"omitempty,max=20" example:"+1234567890"`
}

// LoginRequest represents user login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"test@example.com"`
	Password string `json:"password" binding:"required" example:"securepassword123"`
}

// UpdateUserRequest represents user update request
type UpdateUserRequest struct {
	FirstName *string `json:"first_name,omitempty" binding:"omitempty,max=50" example:"John"`
	LastName  *string `json:"last_name,omitempty" binding:"omitempty,max=50" example:"Doe"`
	Avatar    *string `json:"avatar,omitempty" binding:"omitempty,max=255" example:"https://example.com/avatar.jpg"`
	Phone     *string `json:"phone,omitempty" binding:"omitempty,max=20" example:"+1234567890"`
}

// ChangePasswordRequest represents password change request
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required" example:"oldpassword123"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=50" example:"newpassword123"`
}

// UserListRequest represents user list request with pagination and filters
type UserListRequest struct {
	Page     int              `form:"page,default=1" binding:"min=1" example:"1"`
	Size     int              `form:"size,default=10" binding:"min=1,max=100" example:"10"`
	Sort     string           `form:"sort,default=created_at" example:"created_at"`
	Order    string           `form:"order,default=desc" binding:"oneof=asc desc" example:"desc"`
	Search   string           `form:"search" example:"john"`
	Status   model.UserStatus `form:"status" example:"active"`
	IsActive *bool            `form:"is_active" example:"true"`
}

// RegisterResponse represents user registration response
type RegisterResponse struct {
	User    *model.PublicUser `json:"user"`
	Token   string            `json:"token"`
	Message string            `json:"message"`
}

// LoginResponse represents user login response
type LoginResponse struct {
	User    *model.PublicUser `json:"user"`
	Token   string            `json:"token"`
	Message string            `json:"message"`
}

// UserResponse represents single user response
type UserResponse struct {
	User    *model.PublicUser `json:"user"`
	Message string            `json:"message"`
}

// UserListResponse represents user list response
type UserListResponse struct {
	Users      []*model.PublicUser `json:"users"`
	Pagination *PaginationResponse `json:"pagination"`
	Message    string              `json:"message"`
}

// PaginationResponse represents pagination information
type PaginationResponse struct {
	Page       int   `json:"page"`
	Size       int   `json:"size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// ErrorResponse represents error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

// SuccessResponse represents success response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status    string            `json:"status"`
	Version   string            `json:"version"`
	Timestamp string            `json:"timestamp"`
	Checks    map[string]string `json:"checks"`
}
