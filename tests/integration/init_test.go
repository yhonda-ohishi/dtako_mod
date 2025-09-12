package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
)

// TestMain runs before all tests in this package
func TestMain(m *testing.M) {
	// Load .env file from project root
	projectRoot := filepath.Join("..", "..")
	envPath := filepath.Join(projectRoot, ".env")
	
	if err := godotenv.Load(envPath); err != nil {
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
	}
	
	// Run tests
	code := m.Run()
	
	// Exit with test result code
	os.Exit(code)
}