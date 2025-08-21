package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"analytics-service/internal/analytics"
	"analytics-service/internal/auth"
)

type AnalyticsHandler struct {
	analyticsService *analytics.Service
	authService      *auth.Service
}

func NewAnalyticsHandler(analyticsService *analytics.Service, authService *auth.Service) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
		authService:      authService,
	}
}

func (h *AnalyticsHandler) GetItemAnalytics(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	var req analytics.ItemAnalyticsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Token == "" || req.StartDate == "" || req.FinishDate == "" {
		http.Error(w, "Token, StartDate and FinishDate are required", http.StatusBadRequest)
		return
	}

	validateResponse, err := h.authService.ValidateToken(req.Token)
	if err != nil {
		http.Error(w, "Failed to validate token", http.StatusInternalServerError)
		return
	}

	if !validateResponse.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	response, err := h.analyticsService.GetItemAnalytics(&req)
	if err != nil {
		http.Error(w, "Failed to process analytics", http.StatusInternalServerError)
		return
	}

	processingTime := time.Since(startTime)
	log.Printf("Analytics request processed in %v", processingTime)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
