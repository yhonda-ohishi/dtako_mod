package main

import (
	"log"
	"time"

	"github.com/yhonda-ohishi/dtako_mod/repositories"
)

func main() {
	log.Println("🔍 Testing DTako Events GetByID...")

	// Create repository instance
	repo := repositories.NewDtakoEventsRepository()

	// Test with the ID we got from previous test
	testID := "21103756"

	log.Printf("🔍 Testing GetByID with id=%s", testID)

	// Test the GetByID method
	start := time.Now()
	result, err := repo.GetByID(testID)
	elapsed := time.Since(start)

	if err != nil {
		log.Printf("❌ ERROR: %v (took %v)", err, elapsed)
	} else {
		log.Printf("✅ SUCCESS: Got result (took %v)", elapsed)
		log.Printf("📄 Result: ID=%s, UnkoNo=%s, EventType=%s, EventDate=%v",
			result.ID, result.UnkoNo, result.EventType, result.EventDate)
		if result.Description != "" {
			log.Printf("📝 Description: %s", result.Description)
		}
		if result.Latitude != nil && result.Longitude != nil {
			log.Printf("📍 Location: Lat=%f, Lng=%f", *result.Latitude, *result.Longitude)
		}
	}

	log.Println("🔍 Test completed")
}