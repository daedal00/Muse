package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	Port        string
	Environment string

	// Spotify
	SpotifyClientID     string
	SpotifyClientSecret string
	SpotifyRedirectURL  string

	// Database
	DatabaseURL string
	DBHost      string
	DBPort      int
	DBName      string
	DBUser      string
	DBPassword  string
	DBSSLMode   string

	// Redis
	RedisURL      string
	RedisHost     string
	RedisPort     int
	RedisPassword string
	RedisDB       int

	// JWT
	JWTSecret string
}

func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load(".env")

	// Handle Redis URL configuration - check multiple environment variables
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		// Check for REDIS_ADDR as alternative (used in some CI environments)
		if redisAddr := os.Getenv("REDIS_ADDR"); redisAddr != "" {
			redisDB := getEnvAsInt("REDIS_DB", 0)
			redisPassword := os.Getenv("REDIS_PASSWORD")

			// Construct Redis URL with optional password and database
			if redisPassword != "" {
				if redisDB != 0 {
					redisURL = fmt.Sprintf("redis://:%s@%s/%d", redisPassword, redisAddr, redisDB)
				} else {
					redisURL = fmt.Sprintf("redis://:%s@%s", redisPassword, redisAddr)
				}
			} else {
				if redisDB != 0 {
					redisURL = fmt.Sprintf("redis://%s/%d", redisAddr, redisDB)
				} else {
					redisURL = "redis://" + redisAddr
				}
			}
		} else {
			redisURL = "redis://localhost:6379"
		}
	}

	config := &Config{
		Port:        getEnv("PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),

		SpotifyClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		SpotifyClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		SpotifyRedirectURL:  getEnv("SPOTIFY_REDIRECT_URL", "https://127.0.0.1:8080/auth/spotify/callback"),

		DatabaseURL: os.Getenv("DATABASE_URL"),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnvAsInt("DB_PORT", 5432),
		DBName:      getEnv("DB_NAME", "muse"),
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPassword:  os.Getenv("DB_PASSWORD"),
		DBSSLMode:   getEnv("DB_SSL_MODE", "prefer"),

		RedisURL:      redisURL,
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnvAsInt("REDIS_PORT", 6379),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       getEnvAsInt("REDIS_DB", 0),

		JWTSecret: getEnv("JWT_SECRET", "your-fallback-secret-key"),
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

func (c *Config) Validate() error {
	if c.SpotifyClientID == "" {
		return fmt.Errorf("SPOTIFY_CLIENT_ID is required")
	}
	if c.SpotifyClientSecret == "" {
		return fmt.Errorf("SPOTIFY_CLIENT_SECRET is required")
	}
	if c.DatabaseURL == "" && c.DBPassword == "" {
		return fmt.Errorf("either DATABASE_URL or DB_PASSWORD must be provided")
	}
	return nil
}

func (c *Config) GetDatabaseDSN() string {
	if c.DatabaseURL != "" {
		return c.DatabaseURL
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode)
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
			return intValue
		}
	}
	return defaultValue
}
