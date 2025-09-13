package repositories

import (
	"database/sql"
	"log"
	"os"
	"sync"

	"github.com/yhonda-ohishi/dtako_mod/config"
)

var (
	db    *sql.DB
	once  sync.Once
	dbErr error
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

// GetProductionDB returns the production database connection (uses same as local for now)
func GetProductionDB() (*sql.DB, error) {
	return GetDB()
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