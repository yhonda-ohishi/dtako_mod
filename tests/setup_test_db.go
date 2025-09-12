package tests

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// SetupTestDB creates test databases and loads test data
func SetupTestDB() error {
	// Load .env file
	envPath := filepath.Join(".env")
	if err := godotenv.Load(envPath); err != nil {
		fmt.Printf("Warning: Could not load .env file: %v\n", err)
	}

	// Setup production test database
	if err := setupDatabase("PROD"); err != nil {
		return fmt.Errorf("failed to setup production test database: %v", err)
	}

	// Setup local test database
	if err := setupDatabase("LOCAL"); err != nil {
		return fmt.Errorf("failed to setup local test database: %v", err)
	}

	return nil
}

func setupDatabase(prefix string) error {
	host := os.Getenv(prefix + "_DB_HOST")
	port := os.Getenv(prefix + "_DB_PORT")
	user := os.Getenv(prefix + "_DB_USER")
	password := os.Getenv(prefix + "_DB_PASSWORD")
	dbName := os.Getenv(prefix + "_DB_NAME")

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "3307"
	}
	if user == "" {
		user = "root"
	}
	if password == "" {
		password = "kikuraku"
	}
	if dbName == "" {
		if prefix == "PROD" {
			dbName = "dtako_test_prod"
		} else {
			dbName = "dtako_local"
		}
	}

	// Connect to MySQL without database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?parseTime=true&charset=utf8mb4",
		user, password, host, port)
	
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to MySQL: %v", err)
	}
	defer db.Close()

	// Create database if not exists
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName))
	if err != nil {
		return fmt.Errorf("failed to create database %s: %v", dbName, err)
	}

	// Connect to the specific database
	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		user, password, host, port, dbName)
	
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database %s: %v", dbName, err)
	}
	defer db.Close()

	// Load and execute schema
	schemaPath := filepath.Join("tests", "testdata", "schema.sql")
	if err := executeSQLFile(db, schemaPath); err != nil {
		return fmt.Errorf("failed to execute schema for %s: %v", dbName, err)
	}

	// For production test database, also load test data
	if prefix == "PROD" {
		dataPath := filepath.Join("tests", "testdata", "test_data.sql")
		if err := executeSQLFile(db, dataPath); err != nil {
			return fmt.Errorf("failed to load test data for %s: %v", dbName, err)
		}
	}

	fmt.Printf("Successfully setup %s database: %s\n", prefix, dbName)
	return nil
}

func executeSQLFile(db *sql.DB, filepath string) error {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read SQL file %s: %v", filepath, err)
	}

	// Split by semicolon and execute each statement
	statements := strings.Split(string(content), ";")
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" || strings.HasPrefix(stmt, "--") {
			continue
		}
		
		if _, err := db.Exec(stmt); err != nil {
			return fmt.Errorf("failed to execute statement: %v\nStatement: %s", err, stmt)
		}
	}

	return nil
}