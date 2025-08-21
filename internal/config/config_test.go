package config

import (
	"os"
	"testing"
)

func TestConfig_New(t *testing.T) {

	os.Unsetenv("PORT")
	os.Unsetenv("AUTH_SECRET_KEY")
	os.Unsetenv("WORKERS")
	
	cfg := New()
	

	if cfg.Port != "8080" {
		t.Fatalf("Expected default port 8080, got %s", cfg.Port)
	}
	
	if cfg.SecretKey != "secret" {
		t.Fatalf("Expected default secret key 'secret', got %s", cfg.SecretKey)
	}
	
	if cfg.Workers != 4 {
		t.Fatalf("Expected default workers 4, got %d", cfg.Workers)
	}
}

func TestConfig_New_WithEnvironmentVariables(t *testing.T) {

	os.Setenv("PORT", "9090")
	os.Setenv("AUTH_SECRET_KEY", "test-secret")
	os.Setenv("WORKERS", "8")
	
	cfg := New()
	

	if cfg.Port != "9090" {
		t.Fatalf("Expected port 9090, got %s", cfg.Port)
	}
	
	if cfg.SecretKey != "test-secret" {
		t.Fatalf("Expected secret key 'test-secret', got %s", cfg.SecretKey)
	}
	
	if cfg.Workers != 8 {
		t.Fatalf("Expected workers 8, got %d", cfg.Workers)
	}
	

	os.Unsetenv("PORT")
	os.Unsetenv("AUTH_SECRET_KEY")
	os.Unsetenv("WORKERS")
}

func TestConfig_New_InvalidWorkers(t *testing.T) {

	os.Setenv("WORKERS", "invalid")
	
	cfg := New()
	

	if cfg.Workers != 4 {
		t.Fatalf("Expected default workers 4 for invalid value, got %d", cfg.Workers)
	}
	

	os.Unsetenv("WORKERS")
}

func TestConfig_New_EmptyEnvironmentVariables(t *testing.T) {

	os.Setenv("PORT", "")
	os.Setenv("AUTH_SECRET_KEY", "")
	os.Setenv("WORKERS", "")
	
	cfg := New()
	

	if cfg.Port != "8080" {
		t.Fatalf("Expected default port 8080 for empty value, got %s", cfg.Port)
	}
	
	if cfg.SecretKey != "secret" {
		t.Fatalf("Expected default secret key 'secret' for empty value, got %s", cfg.SecretKey)
	}
	
	if cfg.Workers != 4 {
		t.Fatalf("Expected default workers 4 for empty value, got %d", cfg.Workers)
	}
	

	os.Unsetenv("PORT")
	os.Unsetenv("AUTH_SECRET_KEY")
	os.Unsetenv("WORKERS")
}

func TestConfig_New_NegativeWorkers(t *testing.T) {

	os.Setenv("WORKERS", "-5")
	
	cfg := New()
	

	if cfg.Workers != 4 {
		t.Fatalf("Expected default workers 4 for negative value, got %d", cfg.Workers)
	}
	

	os.Unsetenv("WORKERS")
}

func TestConfig_New_ZeroWorkers(t *testing.T) {

	os.Setenv("WORKERS", "0")
	
	cfg := New()
	

	if cfg.Workers != 4 {
		t.Fatalf("Expected default workers 4 for zero value, got %d", cfg.Workers)
	}
	

	os.Unsetenv("WORKERS")
}

func TestConfig_New_LargeWorkers(t *testing.T) {

	os.Setenv("WORKERS", "100")
	
	cfg := New()
	

	if cfg.Workers != 100 {
		t.Fatalf("Expected workers 100, got %d", cfg.Workers)
	}
	

	os.Unsetenv("WORKERS")
}
