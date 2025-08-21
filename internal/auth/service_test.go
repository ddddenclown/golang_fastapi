package auth

import (
	"testing"

	"analytics-service/internal/userdb"
)

func TestService_GenerateToken(t *testing.T) {
	tokenStore := userdb.NewTokenStore()
	service := NewService("test-secret", tokenStore)

	req := &AuthRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	response, err := service.GenerateToken(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if response.Token == "" {
		t.Fatal("Expected non-empty token")
	}


	userID, exists := tokenStore.ValidateToken(response.Token)
	if !exists {
		t.Fatal("Expected token to be in store")
	}

	if userID != 1 {
		t.Fatalf("Expected user ID 1, got %d", userID)
	}
}

func TestService_ValidateToken(t *testing.T) {
	tokenStore := userdb.NewTokenStore()
	service := NewService("test-secret", tokenStore)


	testToken := "test-token-123"
	tokenStore.AddToken(testToken, 1)


	response, err := service.ValidateToken(testToken)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !response.Valid {
		t.Fatal("Expected token to be valid")
	}


	response, err = service.ValidateToken("invalid-token")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if response.Valid {
		t.Fatal("Expected token to be invalid")
	}
}

func TestService_ValidateTokenWithWrongUserID(t *testing.T) {
	tokenStore := userdb.NewTokenStore()
	service := NewService("test-secret", tokenStore)


	testToken := "test-token-123"
	tokenStore.AddToken(testToken, 2)


	response, err := service.ValidateToken(testToken)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if response.Valid {
		t.Fatal("Expected token to be invalid due to wrong user ID")
	}
}
