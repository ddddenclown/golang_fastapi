package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port      string
	SecretKey string
	Workers   int
}

func New() *Config {
	workers := getEnvAsInt("WORKERS", 4)
	
	return &Config{
		Port:      getEnv("PORT", "8080"),
		SecretKey: getEnv("AUTH_SECRET_KEY", "secret"),
		Workers:   workers,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			if intValue > 0 {
				return intValue
			}
		}
	}
	return defaultValue
}
