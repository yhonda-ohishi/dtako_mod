package main

import (
	"log"
	"time"

	"github.com/yhonda-ohishi/dtako_mod/repositories"
)

func main() {
	log.Println("🔍 Testing DTako Events with Production DB...")

	// Create repository instance
	repo := repositories.NewDtakoEventsRepository()

	// Test dates
	from, _ := time.Parse("2006-01-02", "2025-09-13")
	to, _ := time.Parse("2006-01-02", "2025-09-14")

	log.Printf("🔍 Testing GetByDateRange with from=%v, to=%v", from, to)

	// Test the actual method
	start := time.Now()
	results, err := repo.GetByDateRange(from, to, "", "")
	elapsed := time.Since(start)

	if err != nil {
		log.Printf("❌ ERROR: %v (took %v)", err, elapsed)
	} else {
		log.Printf("✅ SUCCESS: Got %d results (took %v)", len(results), elapsed)

		// Show first result if any
		if len(results) > 0 {
			first := results[0]
			log.Printf("📄 First result: ID=%s, UnkoNo=%s, EventType=%s, EventDate=%v",
				first.ID, first.UnkoNo, first.EventType, first.EventDate)
		}
	}

	log.Println("🔍 Test completed")
}