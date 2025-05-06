package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// Config holds all configuration for the application
type Config struct {
	Host                  string
	Port                  string
	JWTSecret             string
	JWTExpirationMinutes  int
	MicrosoftClientID     string
	MicrosoftClientSecret string
	MicrosoftRedirectURI  string
	MicrosoftTenantID     string
	AppURL                string

	// Database configuration
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	// Load .env file if it exists
	_ = godotenv.Load()

	// Set default values
	config := &Config{
		Host:                  getEnv("HOST", ""),
		Port:                  getEnv("PORT", "8080"),
		JWTSecret:             getEnv("JWT_SECRET", "your-secret-key"),
		JWTExpirationMinutes:  60, // 1 hour
		MicrosoftClientID:     getEnv("MICROSOFT_CLIENT_ID", ""),
		MicrosoftClientSecret: getEnv("MICROSOFT_CLIENT_SECRET", ""),
		MicrosoftRedirectURI:  getEnv("MICROSOFT_REDIRECT_URI", "http://localhost:8080/auth/microsoft/callback"),
		MicrosoftTenantID:     getEnv("MICROSOFT_TENANT_ID", "common"),
		AppURL:                getEnv("APP_URL", "http://localhost:3000"),

		// Database configuration
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "go_azure"),
	}

	// Log configuration
	logrus.Info("Configuration loaded")

	return config
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
