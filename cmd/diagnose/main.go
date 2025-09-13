package main

import (
	"fmt"
	"log"
	"os"

	"github.com/yhonda-ohishi/dtako_mod/config"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	fmt.Println("=== DTako Database Connection Diagnostic ===")

	// 環境変数の確認
	fmt.Println("\n[Environment Variables]")
	fmt.Printf("DB_HOST: %s\n", os.Getenv("DB_HOST"))
	fmt.Printf("DB_PORT: %s\n", os.Getenv("DB_PORT"))
	fmt.Printf("DB_USER: %s\n", os.Getenv("DB_USER"))
	fmt.Printf("DB_PASSWORD: %s\n", maskPassword(os.Getenv("DB_PASSWORD")))
	fmt.Printf("DB_NAME: %s\n", os.Getenv("DB_NAME"))

	// Fallback環境変数の確認
	fmt.Println("\n[Fallback Environment Variables]")
	fmt.Printf("LOCAL_DB_HOST: %s\n", os.Getenv("LOCAL_DB_HOST"))
	fmt.Printf("LOCAL_DB_PORT: %s\n", os.Getenv("LOCAL_DB_PORT"))
	fmt.Printf("LOCAL_DB_USER: %s\n", os.Getenv("LOCAL_DB_USER"))
	fmt.Printf("LOCAL_DB_PASSWORD: %s\n", maskPassword(os.Getenv("LOCAL_DB_PASSWORD")))
	fmt.Printf("LOCAL_DB_NAME: %s\n", os.Getenv("LOCAL_DB_NAME"))

	// 接続テスト
	fmt.Println("\n[Connection Test]")
	cfg := config.GetDatabaseConfig()
	fmt.Printf("Using configuration: %s:%s/%s\n", cfg.Host, cfg.Port, cfg.Database)

	db, err := cfg.Connect()
	if err != nil {
		log.Fatalf("❌ Connection failed: %v", err)
	}
	defer db.Close()

	fmt.Println("✅ Database connection successful!")

	// テーブル確認
	fmt.Println("\n[Tables Check]")
	tables := []string{"dtako_rows", "dtako_events", "dtako_ferry_rows"}
	for _, table := range tables {
		var exists string
		err := db.QueryRow("SHOW TABLES LIKE ?", table).Scan(&exists)
		if err != nil {
			fmt.Printf("❌ Table %s not found\n", table)
		} else {
			fmt.Printf("✅ Table %s exists\n", table)
		}
	}

	// レコード数の確認
	fmt.Println("\n[Record Count]")
	for _, table := range tables {
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&count)
		if err != nil {
			fmt.Printf("   %s: Error counting records: %v\n", table, err)
		} else {
			fmt.Printf("   %s: %d records\n", table, count)
		}
	}
}

func maskPassword(password string) string {
	if len(password) == 0 {
		return "(empty)"
	}
	return "***"
}