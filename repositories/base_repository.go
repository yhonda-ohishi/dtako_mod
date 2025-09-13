package repositories

import (
	"database/sql"
	"log"
	"os"
	"sync"

	"github.com/yhonda-ohishi/dtako_mod/config"
)

var (
	db      *sql.DB
	prodDB  *sql.DB
	once    sync.Once
	onceProd sync.Once
	dbErr   error
	prodErr error
)

// GetDB returns a singleton database connection
func GetDB() (*sql.DB, error) {
	once.Do(func() {
		cfg := config.GetDatabaseConfig()
		db, dbErr = cfg.Connect()
		if dbErr != nil {
			log.Printf("Failed to connect to database: %v", dbErr)
		}
	})
	return db, dbErr
}

// GetLocalDB returns the local database connection (alias for GetDB)
func GetLocalDB() (*sql.DB, error) {
	return GetDB()
}

// GetProductionDB returns the production database connection
func GetProductionDB() (*sql.DB, error) {
	onceProd.Do(func() {
		// Production database configuration from PROD_DB_* env vars
		cfg := &config.DatabaseConfig{
			Host:     getEnvWithDefault("PROD_DB_HOST", "localhost"),
			Port:     getEnvWithDefault("PROD_DB_PORT", "3306"),
			User:     getEnvWithDefault("PROD_DB_USER", "root"),
			Password: getEnvWithDefault("PROD_DB_PASSWORD", ""),
			Database: getEnvWithDefault("PROD_DB_NAME", "dtako_test_prod"),
			Charset:  getEnvWithDefault("PROD_DB_CHARSET", "utf8mb4"),
		}
		prodDB, prodErr = cfg.Connect()
		if prodErr != nil {
			log.Printf("Failed to connect to production database: %v", prodErr)
		}
	})
	return prodDB, prodErr
}

// getEnvWithDefault gets environment variable with default value
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// SetDatabaseConfig allows external configuration (optional)
func SetDatabaseConfig(host, port, user, password, database string) error {
	// 環境変数を設定
	os.Setenv("DB_HOST", host)
	os.Setenv("DB_PORT", port)
	os.Setenv("DB_USER", user)
	os.Setenv("DB_PASSWORD", password)
	os.Setenv("DB_NAME", database)

	// 再接続
	cfg := config.GetDatabaseConfig()
	newDB, err := cfg.Connect()
	if err != nil {
		return err
	}

	// 古い接続をクローズ
	if db != nil {
		db.Close()
	}

	db = newDB
	return nil
}