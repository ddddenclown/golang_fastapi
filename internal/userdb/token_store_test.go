package userdb

import (
	"os"
	"testing"
)

func TestTokenStore_AddAndValidateToken(t *testing.T) {
	store := NewTokenStore()


	token := "test-token"
	userID := 1
	store.AddToken(token, userID)


	retrievedUserID, exists := store.ValidateToken(token)
	if !exists {
		t.Fatal("Expected token to exist")
	}

	if retrievedUserID != userID {
		t.Fatalf("Expected user ID %d, got %d", userID, retrievedUserID)
	}


	_, exists = store.ValidateToken("non-existent")
	if exists {
		t.Fatal("Expected token to not exist")
	}
}

func TestTokenStore_LoadFromFile(t *testing.T) {

	content := "token1 1\ntoken2 2\ntoken3 3"
	tmpFile, err := os.CreateTemp("", "test-tokens")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())


	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()


	store := NewTokenStore()
	if err := store.LoadFromFile(tmpFile.Name()); err != nil {
		t.Fatalf("Failed to load tokens: %v", err)
	}


	testCases := []struct {
		token  string
		userID int
		exists bool
	}{
		{"token1", 1, true},
		{"token2", 2, true},
		{"token3", 3, true},
		{"non-existent", 0, false},
	}

	for _, tc := range testCases {
		userID, exists := store.ValidateToken(tc.token)
		if exists != tc.exists {
			t.Fatalf("Token %s: expected exists=%v, got %v", tc.token, tc.exists, exists)
		}
		if exists && userID != tc.userID {
			t.Fatalf("Token %s: expected user ID %d, got %d", tc.token, tc.userID, userID)
		}
	}
}

func TestTokenStore_GetTokenCount(t *testing.T) {
	store := NewTokenStore()


	if count := store.GetTokenCount(); count != 0 {
		t.Fatalf("Expected 0 tokens, got %d", count)
	}


	store.AddToken("token1", 1)
	store.AddToken("token2", 2)


	if count := store.GetTokenCount(); count != 2 {
		t.Fatalf("Expected 2 tokens, got %d", count)
	}
}
