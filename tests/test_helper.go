package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
)

// SetupTestEnv loads the .env file from the project root
func SetupTestEnv(t *testing.T) {
	t.Helper()
	
	// Find the project root (where .env is located)
	// Go up from tests directory
	projectRoot := filepath.Join("..", "..")
	envPath := filepath.Join(projectRoot, ".env")
	
	// Load .env file
	if err := godotenv.Load(envPath); err != nil {
		t.Logf("Warning: Could not load .env file from %s: %v", envPath, err)
		
		// Set test defaults if .env not found
		os.Setenv("PROD_DB_HOST", "localhost")
		os.Setenv("PROD_DB_PORT", "3307")
		os.Setenv("PROD_DB_USER", "root")
		os.Setenv("PROD_DB_PASSWORD", "kikuraku")
		os.Setenv("PROD_DB_NAME", "dtako_test_prod")
		os.Setenv("PROD_DB_CHARSET", "utf8mb4")
		
		os.Setenv("LOCAL_DB_HOST", "localhost")
		os.Setenv("LOCAL_DB_PORT", "3307")
		os.Setenv("LOCAL_DB_USER", "root")
		os.Setenv("LOCAL_DB_PASSWORD", "kikuraku")
		os.Setenv("LOCAL_DB_NAME", "dtako_local")
		os.Setenv("LOCAL_DB_CHARSET", "utf8mb4")
		
		t.Log("Using default test environment variables")
	} else {
		t.Logf("Loaded .env file from %s", envPath)
	}
}