package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/your-org/user-center/internal/dto"
	"github.com/your-org/user-center/pkg/jwt"
	"go.uber.org/zap"
)

// AuthMiddleware handles JWT authentication
type AuthMiddleware struct {
	jwtManager *jwt.JWT
	logger     *zap.Logger
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(jwtManager *jwt.JWT, logger *zap.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
		logger:     logger,
	}
}

// RequireAuth validates JWT token and sets user claims in context
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			m.logger.Warn("Missing authorization header")
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error:   "Unauthorized",
				Message: "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			m.logger.Warn("Invalid authorization header format")
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error:   "Unauthorized",
				Message: "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		token := parts[1]

		// Validate token
		claims, err := m.jwtManager.ValidateToken(token)
		if err != nil {
			m.logger.Warn("Invalid JWT token", zap.Error(err))
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error:   "Unauthorized",
				Message: "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Set claims in context
		c.Set("claims", claims)
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)

		c.Next()
	}
}

// OptionalAuth validates JWT token if present but doesn't require it
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		token := parts[1]

		// Validate token
		claims, err := m.jwtManager.ValidateToken(token)
		if err != nil {
			m.logger.Debug("Invalid optional JWT token", zap.Error(err))
			c.Next()
			return
		}

		// Set claims in context
		c.Set("claims", claims)
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)

		c.Next()
	}
}

// RequireActiveUser ensures the authenticated user is active
func (m *AuthMiddleware) RequireActiveUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error:   "Unauthorized",
				Message: "Authentication required",
			})
			c.Abort()
			return
		}

		userClaims := claims.(*jwt.Claims)
		if userClaims.Status != "active" {
			m.logger.Warn("Inactive user attempting to access protected resource",
				zap.Uint("user_id", userClaims.UserID),
				zap.String("status", string(userClaims.Status)),
			)
			c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error:   "Forbidden",
				Message: "Account is not active",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// AdminOnly ensures the authenticated user has admin privileges
func (m *AuthMiddleware) AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error:   "Unauthorized",
				Message: "Authentication required",
			})
			c.Abort()
			return
		}

		userClaims := claims.(*jwt.Claims)

		// Note: This is a simple check. In a real application, you would
		// check user roles from the database or include roles in JWT claims
		if userClaims.Email != "admin@example.com" {
			m.logger.Warn("Non-admin user attempting to access admin resource",
				zap.Uint("user_id", userClaims.UserID),
				zap.String("email", userClaims.Email),
			)
			c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error:   "Forbidden",
				Message: "Admin access required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
