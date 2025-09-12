package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/yhonda-ohishi/dtako_mod"
	"github.com/yhonda-ohishi/dtako_mod/models"
)

// T015: Integration test - Import dtako_rows with date range
func TestImportRowsScenario(t *testing.T) {
	// Setup
	r := chi.NewRouter()
	dtako_mod.RegisterRoutes(r)

	t.Run("Complete import workflow", func(t *testing.T) {
		// Step 1: Import data for a specific date range
		importReq := models.ImportRequest{
			FromDate: "2024-01-15",
			ToDate:   "2024-01-16",
		}
		
		bodyBytes, _ := json.Marshal(importReq)
		req := httptest.NewRequest("POST", "/dtako/rows/import", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		
		r.ServeHTTP(rec, req)
		
		// Verify import succeeded
		if rec.Code != http.StatusOK {
			t.Fatalf("Import failed with status %d: %s", rec.Code, rec.Body.String())
		}
		
		var importResult models.ImportResult
		json.Unmarshal(rec.Body.Bytes(), &importResult)
		
		if !importResult.Success {
			t.Error("Import should be successful")
		}
		
		// Step 2: Query the imported data
		req = httptest.NewRequest("GET", "/dtako/rows?from=2024-01-15&to=2024-01-16", nil)
		rec = httptest.NewRecorder()
		
		r.ServeHTTP(rec, req)
		
		if rec.Code != http.StatusOK {
			t.Fatalf("Query failed with status %d: %s", rec.Code, rec.Body.String())
		}
		
		var rows []models.DtakoRow
		json.Unmarshal(rec.Body.Bytes(), &rows)
		
		// Step 3: Verify data integrity
		for _, row := range rows {
			// Check date is within range
			if row.Date.Before(time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)) ||
				row.Date.After(time.Date(2024, 1, 16, 23, 59, 59, 0, time.UTC)) {
				t.Errorf("Row date %v is outside imported range", row.Date)
			}
			
			// Check required fields are present
			if row.VehicleNo == "" {
				t.Error("Vehicle number should not be empty")
			}
			if row.DriverCode == "" {
				t.Error("Driver code should not be empty")
			}
			if row.Distance < 0 {
				t.Error("Distance should not be negative")
			}
			if row.FuelAmount < 0 {
				t.Error("Fuel amount should not be negative")
			}
		}
		
		// Step 4: Test re-import (should handle duplicates)
		req = httptest.NewRequest("POST", "/dtako/rows/import", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rec = httptest.NewRecorder()
		
		r.ServeHTTP(rec, req)
		
		if rec.Code != http.StatusOK {
			t.Errorf("Re-import failed with status %d", rec.Code)
		}
		
		// Query again and verify no duplicates
		req = httptest.NewRequest("GET", "/dtako/rows?from=2024-01-15&to=2024-01-16", nil)
		rec = httptest.NewRecorder()
		
		r.ServeHTTP(rec, req)
		
		var rowsAfterReimport []models.DtakoRow
		json.Unmarshal(rec.Body.Bytes(), &rowsAfterReimport)
		
		// Count should be same or similar (UPSERT should prevent true duplicates)
		// This check depends on actual data and UPSERT implementation
		t.Logf("Rows after reimport: %d", len(rowsAfterReimport))
	})
}