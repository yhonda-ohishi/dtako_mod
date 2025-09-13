package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload" // .envファイルを自動読み込み
)

// GetDatabaseConfig returns database configuration from environment variables
// Uses DB_* environment variables as primary, falls back to LOCAL_DB_* for compatibility
func GetDatabaseConfig() *DatabaseConfig {
	config := &DatabaseConfig{
		Host:     getEnvWithFallback("DB_HOST", "LOCAL_DB_HOST", "localhost"),
		Port:     getEnvWithFallback("DB_PORT", "LOCAL_DB_PORT", "3306"),
		User:     getEnvWithFallback("DB_USER", "LOCAL_DB_USER", "root"),
		Password: getEnvWithFallback("DB_PASSWORD", "LOCAL_DB_PASSWORD", ""),
		Database: getEnvWithFallback("DB_NAME", "LOCAL_DB_NAME", "dtako_local"),
		Charset:  "utf8mb4",
	}

	return config
}

// getEnvWithFallback tries primary key, then fallback key, then default value
func getEnvWithFallback(primary, fallback, defaultValue string) string {
	if value := os.Getenv(primary); value != "" {
		return value
	}
	if value := os.Getenv(fallback); value != "" {
		return value
	}
	return defaultValue
}

// GetDSN returns the database connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.Database)
}

// Connect establishes database connection
func (c *DatabaseConfig) Connect() (*sql.DB, error) {
	dsn := c.GetDSN()

	// デバッグログ（環境変数DEBUGがtrueの場合のみ）
	if os.Getenv("DEBUG") == "true" {
		fmt.Printf("[DEBUG] Connecting to database at %s:%s/%s\n",
			c.Host, c.Port, c.Database)
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 接続テスト
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}