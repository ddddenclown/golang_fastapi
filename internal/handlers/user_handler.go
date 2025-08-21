package handlers

import (
	"encoding/json"
	"net/http"

	"analytics-service/internal/auth"
)

type UserHandler struct {
	authService *auth.Service
}

func NewUserHandler(authService *auth.Service) *UserHandler {
	return &UserHandler{
		authService: authService,
	}
}

func (h *UserHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Token parameter is required", http.StatusBadRequest)
		return
	}

	response, err := h.authService.ValidateToken(token)
	if err != nil {
		http.Error(w, "Failed to validate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
