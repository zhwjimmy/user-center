package jwt

import (
	"testing"
	"time"
)

// MockUser implements the User interface for testing
type MockUser struct {
	ID       uint
	Username string
	Email    string
	Status   string
}

func (m *MockUser) GetID() uint {
	return m.ID
}

func (m *MockUser) GetUsername() string {
	return m.Username
}

func (m *MockUser) GetEmail() string {
	return m.Email
}

func (m *MockUser) GetStatus() string {
	return m.Status
}

func TestJWT_GenerateAndValidateToken(t *testing.T) {
	secret := "test-secret-key"
	issuer := "test-issuer"
	expiry := time.Hour

	jwtManager := NewJWT(secret, issuer, expiry)

	user := &MockUser{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Status:   "active",
	}

	// Generate token
	token, err := jwtManager.GenerateToken(user)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Fatal("Generated token is empty")
	}

	// Validate token
	claims, err := jwtManager.ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	// Verify claims
	if claims.UserID != user.GetID() {
		t.Errorf("Expected UserID %d, got %d", user.GetID(), claims.UserID)
	}

	if claims.Username != user.GetUsername() {
		t.Errorf("Expected Username %s, got %s", user.GetUsername(), claims.Username)
	}

	if claims.Email != user.GetEmail() {
		t.Errorf("Expected Email %s, got %s", user.GetEmail(), claims.Email)
	}

	if claims.Status != UserStatusActive {
		t.Errorf("Expected Status %s, got %s", UserStatusActive, claims.Status)
	}

	if claims.Issuer != issuer {
		t.Errorf("Expected Issuer %s, got %s", issuer, claims.Issuer)
	}
}

func TestJWT_ValidateInvalidToken(t *testing.T) {
	secret := "test-secret-key"
	issuer := "test-issuer"
	expiry := time.Hour

	jwtManager := NewJWT(secret, issuer, expiry)

	// Test with invalid token
	_, err := jwtManager.ValidateToken("invalid-token")
	if err == nil {
		t.Fatal("Expected error for invalid token, got nil")
	}
}
