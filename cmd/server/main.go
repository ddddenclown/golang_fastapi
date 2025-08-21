package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"analytics-service/internal/auth"
	"analytics-service/internal/analytics"
	"analytics-service/internal/config"
	"analytics-service/internal/handlers"
	"analytics-service/internal/userdb"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	cfg := config.New()

	tokenStore := userdb.NewTokenStore()
	if err := tokenStore.LoadFromFile("routes/LogPas.txt"); err != nil {
		log.Printf("Warning: could not load tokens from file: %v", err)
	}

	authService := auth.NewService(cfg.SecretKey, tokenStore)

	analyticsService := analytics.NewService()
	analyticsService.SetWorkers(cfg.Workers)

	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(authService)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService, authService)

	router := mux.NewRouter()
	
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			
			next.ServeHTTP(w, r)
		})
	})

	router.HandleFunc("/auth", authHandler.GenerateToken).Methods("POST")
	router.HandleFunc("/validate", userHandler.ValidateToken).Methods("GET")
	router.HandleFunc("/analytics", analyticsHandler.GetItemAnalytics).Methods("POST")
	
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Analytics Service is running"}`))
	}).Methods("GET")

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Starting server on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
