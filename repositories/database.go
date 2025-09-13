package repositories

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/yhonda-ohishi/dtako_mod/config"
)

var (
	prodDB  *sql.DB
	localDB *sql.DB
	once    sync.Once
	initErr error
)

// InitDatabases initializes database connections
func InitDatabases() error {
	once.Do(func() {
		cfg, err := config.Load()
		if err != nil {
			initErr = fmt.Errorf("failed to load config: %v", err)
			return
		}

		// Production database connection
		prodDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
			cfg.ProductionDB.User,
			cfg.ProductionDB.Password,
			cfg.ProductionDB.Host,
			cfg.ProductionDB.Port,
			cfg.ProductionDB.Database,
			cfg.ProductionDB.Charset,
		)

		prodDB, err = sql.Open("mysql", prodDSN)
		if err != nil {
			initErr = fmt.Errorf("failed to open production DB: %v", err)
			return
		}

		// Local database connection
		localDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
			cfg.LocalDB.User,
			cfg.LocalDB.Password,
			cfg.LocalDB.Host,
			cfg.LocalDB.Port,
			cfg.LocalDB.Database,
			cfg.LocalDB.Charset,
		)

		localDB, err = sql.Open("mysql", localDSN)
		if err != nil {
			initErr = fmt.Errorf("failed to open local DB: %v", err)
			return
		}

		// Test connections
		if err = prodDB.Ping(); err != nil {
			initErr = fmt.Errorf("failed to ping production DB: %v", err)
			return
		}
		if err = localDB.Ping(); err != nil {
			initErr = fmt.Errorf("failed to ping local DB: %v", err)
			return
		}
	})
	return initErr
}

// GetProductionDB returns production database connection
func GetProductionDB() (*sql.DB, error) {
	if prodDB == nil {
		if err := InitDatabases(); err != nil {
			return nil, err
		}
	}
	return prodDB, nil
}

// GetLocalDB returns local database connection
func GetLocalDB() (*sql.DB, error) {
	if localDB == nil {
		if err := InitDatabases(); err != nil {
			return nil, err
		}
	}
	return localDB, nil
}

// InitDB initializes database connection (for compatibility)
func InitDB(dsn string) error {
	// This function is for compatibility with existing code
	// Actual initialization happens in InitDatabases
	return InitDatabases()
}

// CloseDB closes database connections (for compatibility)
func CloseDB() {
	CloseConnections()
}

// CloseConnections closes all database connections
func CloseConnections() {
	if prodDB != nil {
		prodDB.Close()
		prodDB = nil
	}
	if localDB != nil {
		localDB.Close()
		localDB = nil
	}
	// Reset the once flag for re-initialization
	once = sync.Once{}
	initErr = nil
}