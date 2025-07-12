package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// UserStatus represents user status in JWT claims
type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusDeleted   UserStatus = "deleted"
)

// Claims represents JWT claims
type Claims struct {
	UserID   string     `json:"user_id"`
	Username string     `json:"username"`
	Email    string     `json:"email"`
	Status   UserStatus `json:"status"`
	jwt.RegisteredClaims
}

// JWT handles JWT token operations
type JWT struct {
	secret string
	issuer string
	expiry time.Duration
}

// NewJWT creates a new JWT manager
func NewJWT(secret, issuer string, expiry time.Duration) *JWT {
	return &JWT{
		secret: secret,
		issuer: issuer,
		expiry: expiry,
	}
}

// GenerateToken generates a JWT token for a user
// User interface to avoid circular dependency
type User interface {
	GetID() string
	GetUsername() string
	GetEmail() string
	GetStatus() string
}

// GenerateToken generates a JWT token for a user
func (j *JWT) GenerateToken(user User) (string, error) {
	// Convert string status to UserStatus
	var status UserStatus
	switch user.GetStatus() {
	case "active":
		status = UserStatusActive
	case "inactive":
		status = UserStatusInactive
	case "suspended":
		status = UserStatusSuspended
	case "deleted":
		status = UserStatusDeleted
	default:
		status = UserStatusInactive
	}

	claims := &Claims{
		UserID:   user.GetID(),
		Username: user.GetUsername(),
		Email:    user.GetEmail(),
		Status:   status,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    j.issuer,
			Subject:   user.GetID(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}

// ValidateToken validates a JWT token and returns claims
func (j *JWT) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
