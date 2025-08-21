package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"analytics-service/internal/analytics"
	"analytics-service/internal/auth"
	"analytics-service/internal/userdb"
)

func TestAuthHandler_GenerateToken(t *testing.T) {
	tokenStore := userdb.NewTokenStore()
	authService := auth.NewService("test-secret", tokenStore)
	handler := NewAuthHandler(authService)


	reqBody := auth.AuthRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	
	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/auth", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	handler.GenerateToken(w, req)
	
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
	
	var response auth.AuthResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	if response.Token == "" {
		t.Fatal("Expected non-empty token")
	}
}

func TestAuthHandler_GenerateToken_InvalidRequest(t *testing.T) {
	tokenStore := userdb.NewTokenStore()
	authService := auth.NewService("test-secret", tokenStore)
	handler := NewAuthHandler(authService)


	req := httptest.NewRequest("POST", "/auth", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	handler.GenerateToken(w, req)
	
	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400, got %d", w.Code)
	}
}

func TestAuthHandler_GenerateToken_MissingFields(t *testing.T) {
	tokenStore := userdb.NewTokenStore()
	authService := auth.NewService("test-secret", tokenStore)
	handler := NewAuthHandler(authService)


	reqBody := auth.AuthRequest{
		Email: "test@example.com",

	}
	
	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/auth", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	handler.GenerateToken(w, req)
	
	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400, got %d", w.Code)
	}
}

func TestUserHandler_ValidateToken(t *testing.T) {
	tokenStore := userdb.NewTokenStore()
	authService := auth.NewService("test-secret", tokenStore)
	handler := NewUserHandler(authService)


	testToken := "test-token-123"
	tokenStore.AddToken(testToken, 1)


	req := httptest.NewRequest("GET", "/validate?token="+testToken, nil)
	w := httptest.NewRecorder()
	handler.ValidateToken(w, req)
	
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
	
	var response auth.ValidateResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	if !response.Valid {
		t.Fatal("Expected token to be valid")
	}
}

func TestUserHandler_ValidateToken_MissingToken(t *testing.T) {
	tokenStore := userdb.NewTokenStore()
	authService := auth.NewService("test-secret", tokenStore)
	handler := NewUserHandler(authService)


	req := httptest.NewRequest("GET", "/validate", nil)
	w := httptest.NewRecorder()
	handler.ValidateToken(w, req)
	
	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400, got %d", w.Code)
	}
}

func TestUserHandler_ValidateToken_InvalidToken(t *testing.T) {
	tokenStore := userdb.NewTokenStore()
	authService := auth.NewService("test-secret", tokenStore)
	handler := NewUserHandler(authService)


	req := httptest.NewRequest("GET", "/validate?token=invalid-token", nil)
	w := httptest.NewRecorder()
	handler.ValidateToken(w, req)
	
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
	
	var response auth.ValidateResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	if response.Valid {
		t.Fatal("Expected token to be invalid")
	}
}

func TestAnalyticsHandler_GetItemAnalytics(t *testing.T) {
	tokenStore := userdb.NewTokenStore()
	authService := auth.NewService("test-secret", tokenStore)
	analyticsService := analytics.NewService()
	handler := NewAnalyticsHandler(analyticsService, authService)


	testToken := "test-token-123"
	tokenStore.AddToken(testToken, 1)


	reqBody := analytics.ItemAnalyticsRequest{
		Token:      testToken,
		StartDate:  "01.01.2024",
		FinishDate: "31.01.2024",
	}
	
	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/analytics", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	handler.GetItemAnalytics(w, req)
	

	if w.Code == http.StatusBadRequest || w.Code == http.StatusUnauthorized {
		t.Fatalf("Expected status not to be 400/401, got %d", w.Code)
	}
}

func TestAnalyticsHandler_GetItemAnalytics_InvalidToken(t *testing.T) {
	tokenStore := userdb.NewTokenStore()
	authService := auth.NewService("test-secret", tokenStore)
	analyticsService := analytics.NewService()
	handler := NewAnalyticsHandler(analyticsService, authService)


	reqBody := analytics.ItemAnalyticsRequest{
		Token:      "invalid-token",
		StartDate:  "01.01.2024",
		FinishDate: "31.01.2024",
	}
	
	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/analytics", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	handler.GetItemAnalytics(w, req)
	
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("Expected status 401, got %d", w.Code)
	}
}

func TestAnalyticsHandler_GetItemAnalytics_MissingFields(t *testing.T) {
	tokenStore := userdb.NewTokenStore()
	authService := auth.NewService("test-secret", tokenStore)
	analyticsService := analytics.NewService()
	handler := NewAnalyticsHandler(analyticsService, authService)


	reqBody := analytics.ItemAnalyticsRequest{
		Token: "test-token",

	}
	
	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/analytics", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	handler.GetItemAnalytics(w, req)
	
	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400, got %d", w.Code)
	}
}

func TestAnalyticsHandler_GetItemAnalytics_InvalidJSON(t *testing.T) {
	tokenStore := userdb.NewTokenStore()
	authService := auth.NewService("test-secret", tokenStore)
	analyticsService := analytics.NewService()
	handler := NewAnalyticsHandler(analyticsService, authService)


	req := httptest.NewRequest("POST", "/analytics", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	handler.GetItemAnalytics(w, req)
	
	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400, got %d", w.Code)
	}
}
