package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Test with environment variables
	os.Setenv("PORT", "8080")
	os.Setenv("ENVIRONMENT", "test")
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("SPOTIFY_CLIENT_ID", "test-client-id")
	os.Setenv("SPOTIFY_CLIENT_SECRET", "test-client-secret")
	os.Setenv("DATABASE_URL", "postgres://test:test@localhost/test")
	
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("ENVIRONMENT")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("SPOTIFY_CLIENT_ID")
		os.Unsetenv("SPOTIFY_CLIENT_SECRET")
		os.Unsetenv("DATABASE_URL")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.Port != "8080" {
		t.Errorf("Expected port 8080, got %s", cfg.Port)
	}

	if cfg.Environment != "test" {
		t.Errorf("Expected environment test, got %s", cfg.Environment)
	}

	if cfg.JWTSecret != "test-secret" {
		t.Errorf("Expected JWT secret test-secret, got %s", cfg.JWTSecret)
	}
}

func TestConfigDefaults(t *testing.T) {
	// Set required environment variables
	os.Setenv("SPOTIFY_CLIENT_ID", "test-client-id")
	os.Setenv("SPOTIFY_CLIENT_SECRET", "test-client-secret")
	os.Setenv("DATABASE_URL", "postgres://test:test@localhost/test")
	
	defer func() {
		os.Unsetenv("SPOTIFY_CLIENT_ID")
		os.Unsetenv("SPOTIFY_CLIENT_SECRET")
		os.Unsetenv("DATABASE_URL")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config with defaults: %v", err)
	}

	// Test default values
	if cfg.Port != "8080" {
		t.Errorf("Expected default port 8080, got %s", cfg.Port)
	}

	if cfg.Environment != "development" {
		t.Errorf("Expected default environment development, got %s", cfg.Environment)
	}
}

func TestConfigValidation(t *testing.T) {
	// Test missing required fields
	os.Clearenv()
	
	_, err := Load()
	if err == nil {
		t.Error("Expected validation error for missing required fields")
	}
}

func BenchmarkLoadConfig(b *testing.B) {
	// Set up test environment
	os.Setenv("PORT", "8080")
	os.Setenv("ENVIRONMENT", "test")
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("SPOTIFY_CLIENT_ID", "test-client-id")
	os.Setenv("SPOTIFY_CLIENT_SECRET", "test-client-secret")
	os.Setenv("DATABASE_URL", "postgres://test:test@localhost/test")
	
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("ENVIRONMENT")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("SPOTIFY_CLIENT_ID")
		os.Unsetenv("SPOTIFY_CLIENT_SECRET")
		os.Unsetenv("DATABASE_URL")
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Load()
		if err != nil {
			b.Fatalf("Failed to load config: %v", err)
		}
	}
} 