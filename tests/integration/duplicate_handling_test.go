package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/yhonda-ohishi/dtako_mod"
	"github.com/yhonda-ohishi/dtako_mod/models"
)

// T019: Integration test - Prevent duplicate imports (UPSERT)
func TestDuplicateHandling(t *testing.T) {
	// Setup
	r := chi.NewRouter()
	dtako_mod.RegisterRoutes(r)

	t.Run("UPSERT prevents duplicates on re-import", func(t *testing.T) {
		dateRange := models.ImportRequest{
			FromDate: "2025-01-15",
			ToDate:   "2025-01-15", // Single day for easier testing
		}
		
		// Step 1: First import
		bodyBytes, _ := json.Marshal(dateRange)
		req := httptest.NewRequest("POST", "/dtako/rows/import", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		
		r.ServeHTTP(rec, req)
		
		if rec.Code != http.StatusOK {
			t.Fatalf("First import failed with status %d: %s", rec.Code, rec.Body.String())
		}
		
		var firstImport models.ImportResult
		json.Unmarshal(rec.Body.Bytes(), &firstImport)
		firstCount := firstImport.ImportedRows
		
		// Step 2: Query to get initial data
		req = httptest.NewRequest("GET", "/dtako/rows?from=2025-01-15&to=2025-01-15", nil)
		rec = httptest.NewRecorder()
		
		r.ServeHTTP(rec, req)
		
		var initialRows []models.DtakoRow
		json.Unmarshal(rec.Body.Bytes(), &initialRows)
		
		// Store IDs for comparison
		idMap := make(map[string]models.DtakoRow)
		for _, row := range initialRows {
			idMap[row.ID] = row
		}
		
		// Step 3: Re-import same date range
		req = httptest.NewRequest("POST", "/dtako/rows/import", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rec = httptest.NewRecorder()
		
		r.ServeHTTP(rec, req)
		
		if rec.Code != http.StatusOK {
			t.Fatalf("Re-import failed with status %d: %s", rec.Code, rec.Body.String())
		}
		
		var secondImport models.ImportResult
		json.Unmarshal(rec.Body.Bytes(), &secondImport)
		
		// Step 4: Query again
		req = httptest.NewRequest("GET", "/dtako/rows?from=2025-01-15&to=2025-01-15", nil)
		rec = httptest.NewRecorder()
		
		r.ServeHTTP(rec, req)
		
		var rowsAfterReimport []models.DtakoRow
		json.Unmarshal(rec.Body.Bytes(), &rowsAfterReimport)
		
		// Verify no duplicates (same count)
		if len(rowsAfterReimport) != len(initialRows) {
			t.Errorf("Expected same row count after re-import. Initial: %d, After: %d",
				len(initialRows), len(rowsAfterReimport))
		}
		
		// Verify IDs match (UPSERT should update, not create new)
		for _, row := range rowsAfterReimport {
			if _, exists := idMap[row.ID]; !exists {
				t.Errorf("Unexpected new ID after re-import: %s", row.ID)
			}
		}
		
		t.Logf("First import: %d rows, Re-import handled correctly with %d total rows",
			firstCount, len(rowsAfterReimport))
	})
	
	t.Run("UPSERT updates modified records", func(t *testing.T) {
		// This test simulates updating records in production and re-importing
		dateRange := models.ImportRequest{
			FromDate: "2025-01-20",
			ToDate:   "2025-01-20",
		}
		
		// Import initial data
		bodyBytes, _ := json.Marshal(dateRange)
		req := httptest.NewRequest("POST", "/dtako/events/import", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		
		r.ServeHTTP(rec, req)
		
		if rec.Code != http.StatusOK {
			t.Fatalf("Initial import failed: %s", rec.Body.String())
		}
		
		// Query initial data
		req = httptest.NewRequest("GET", "/dtako/events?from=2025-01-20&to=2025-01-20", nil)
		rec = httptest.NewRecorder()
		
		r.ServeHTTP(rec, req)
		
		var initialEvents []models.DtakoEvent
		json.Unmarshal(rec.Body.Bytes(), &initialEvents)
		initialCount := len(initialEvents)
		
		// Re-import (simulating updated data in production)
		req = httptest.NewRequest("POST", "/dtako/events/import", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rec = httptest.NewRecorder()
		
		r.ServeHTTP(rec, req)
		
		if rec.Code != http.StatusOK {
			t.Fatalf("Re-import failed: %s", rec.Body.String())
		}
		
		// Query after re-import
		req = httptest.NewRequest("GET", "/dtako/events?from=2025-01-20&to=2025-01-20", nil)
		rec = httptest.NewRecorder()
		
		r.ServeHTTP(rec, req)
		
		var eventsAfterUpdate []models.DtakoEvent
		json.Unmarshal(rec.Body.Bytes(), &eventsAfterUpdate)
		
		// Count should remain same (updates, not duplicates)
		if len(eventsAfterUpdate) != initialCount {
			t.Errorf("Event count changed after re-import. Initial: %d, After: %d",
				initialCount, len(eventsAfterUpdate))
		}
		
		t.Logf("UPSERT correctly handled %d events", len(eventsAfterUpdate))
	})
	
	t.Run("Multiple parallel imports handle correctly", func(t *testing.T) {
		// Test that multiple simultaneous imports don't create duplicates
		dateRange := models.ImportRequest{
			FromDate: "2025-01-25",
			ToDate:   "2025-01-25",
		}
		
		bodyBytes, _ := json.Marshal(dateRange)
		
		// Launch multiple imports in parallel
		done := make(chan bool, 3)
		
		for i := 0; i < 3; i++ {
			go func(id int) {
				req := httptest.NewRequest("POST", "/dtako/ferry/import", bytes.NewReader(bodyBytes))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				
				r.ServeHTTP(rec, req)
				
				if rec.Code != http.StatusOK {
					t.Logf("Parallel import %d returned status %d", id, rec.Code)
				}
				done <- true
			}(i)
		}
		
		// Wait for all imports to complete
		for i := 0; i < 3; i++ {
			<-done
		}
		
		// Query final data
		req := httptest.NewRequest("GET", "/dtako/ferry?from=2025-01-25&to=2025-01-25", nil)
		rec := httptest.NewRecorder()
		
		r.ServeHTTP(rec, req)
		
		var finalRecords []models.DtakoFerry
		json.Unmarshal(rec.Body.Bytes(), &finalRecords)
		
		// Check for duplicate IDs
		idSet := make(map[string]bool)
		for _, record := range finalRecords {
			if idSet[record.ID] {
				t.Errorf("Duplicate ID found: %s", record.ID)
			}
			idSet[record.ID] = true
		}
		
		t.Logf("Parallel imports resulted in %d unique records", len(finalRecords))
	})
}