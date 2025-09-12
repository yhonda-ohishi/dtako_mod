package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds database configuration
type Config struct {
	ProductionDB DatabaseConfig
	LocalDB      DatabaseConfig
}

// DatabaseConfig holds individual database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	Charset  string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	godotenv.Load()

	cfg := &Config{
		ProductionDB: DatabaseConfig{
			Host:     getEnv("PROD_DB_HOST", "localhost"),
			Port:     getEnv("PROD_DB_PORT", "3306"),
			User:     getEnv("PROD_DB_USER", "root"),
			Password: getEnv("PROD_DB_PASSWORD", ""),
			Database: getEnv("PROD_DB_NAME", "production"),
			Charset:  getEnv("PROD_DB_CHARSET", "utf8mb4"),
		},
		LocalDB: DatabaseConfig{
			Host:     getEnv("LOCAL_DB_HOST", "localhost"),
			Port:     getEnv("LOCAL_DB_PORT", "3306"),
			User:     getEnv("LOCAL_DB_USER", "root"),
			Password: getEnv("LOCAL_DB_PASSWORD", ""),
			Database: getEnv("LOCAL_DB_NAME", "dtako_local"),
			Charset:  getEnv("LOCAL_DB_CHARSET", "utf8mb4"),
		},
	}

	return cfg, nil
}

// getEnv gets environment variable with default fallback
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}