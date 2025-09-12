package main

import (
	"fmt"
	"log"
	"os"

	"github.com/yhonda-ohishi/dtako_mod/tests"
)

func main() {
	fmt.Println("Setting up test databases...")
	
	if err := tests.SetupTestDB(); err != nil {
		log.Fatalf("Failed to setup test databases: %v", err)
	}

	fmt.Println("Test databases setup completed successfully!")
	os.Exit(0)
}