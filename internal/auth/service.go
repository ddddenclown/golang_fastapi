package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"analytics-service/internal/userdb"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	secretKey  string
	tokenStore *userdb.TokenStore
}

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type ValidateResponse struct {
	Valid bool `json:"valid"`
}

func NewService(secretKey string, tokenStore *userdb.TokenStore) *Service {
	return &Service{
		secretKey:  secretKey,
		tokenStore: tokenStore,
	}
}

func (s *Service) GenerateToken(req *AuthRequest) (*AuthResponse, error) {
	token, err := s.createJWTToken(req.Email, req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to create token: %w", err)
	}

	s.tokenStore.AddToken(token, 1)

	return &AuthResponse{Token: token}, nil
}

func (s *Service) ValidateToken(token string) (*ValidateResponse, error) {
	userID, exists := s.tokenStore.ValidateToken(token)
	if !exists {
		return &ValidateResponse{Valid: false}, nil
	}

	if userID != 1 {
		return &ValidateResponse{Valid: false}, nil
	}

	return &ValidateResponse{Valid: true}, nil
}

func (s *Service) createJWTToken(email, password string) (string, error) {
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}

	payload := map[string]string{
		"email":    email,
		"password": password,
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	headerB64 := base64URLEncode(headerJSON)

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	payloadB64 := base64URLEncode(payloadJSON)

	message := headerB64 + "." + payloadB64

	h := hmac.New(sha256.New, []byte(s.secretKey))
	h.Write([]byte(message))
	signature := h.Sum(nil)
	signatureB64 := base64URLEncode(signature)

	return message + "." + signatureB64, nil
}

func base64URLEncode(data []byte) string {
	encoded := base64.StdEncoding.EncodeToString(data)
	encoded = strings.ReplaceAll(encoded, "+", "-")
	encoded = strings.ReplaceAll(encoded, "/", "_")
	encoded = strings.TrimRight(encoded, "=")
	return encoded
}

func (s *Service) createStandardJWTToken(email, password string) (string, error) {
	claims := jwt.MapClaims{
		"email":    email,
		"password": password,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}
