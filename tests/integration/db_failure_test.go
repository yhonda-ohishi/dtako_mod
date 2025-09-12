package integration

import (
	"os"
	"testing"
	"net/http"
	"net/http/httptest"
	"bytes"
	"encoding/json"

	"github.com/go-chi/chi/v5"
	"github.com/yhonda-ohishi/dtako_mod"
	"github.com/yhonda-ohishi/dtako_mod/models"
)

// T018: Integration test - Handle production DB failure gracefully
func TestProductionDBFailure(t *testing.T) {
	// Setup
	r := chi.NewRouter()
	dtako_mod.RegisterRoutes(r)

	t.Run("Handle production database unavailable", func(t *testing.T) {
		// Save original env vars
		origHost := os.Getenv("PROD_DB_HOST")
		origPort := os.Getenv("PROD_DB_PORT")
		
		// Set invalid production database connection
		os.Setenv("PROD_DB_HOST", "invalid-host")
		os.Setenv("PROD_DB_PORT", "99999")
		
		// Attempt to import - should fail gracefully
		importReq := models.ImportRequest{
			FromDate: "2025-01-01",
			ToDate:   "2025-01-31",
		}
		
		bodyBytes, _ := json.Marshal(importReq)
		req := httptest.NewRequest("POST", "/dtako/rows/import", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		
		r.ServeHTTP(rec, req)
		
		// Should return 200 with success:false or 500 error
		if rec.Code == http.StatusOK {
			// If 200, check for success:false
			var result models.ImportResult
			if err := json.Unmarshal(rec.Body.Bytes(), &result); err != nil {
				t.Errorf("Failed to parse response: %v", err)
			} else if result.Success {
				t.Error("Expected success:false when production DB is unavailable")
			}
		} else if rec.Code != http.StatusInternalServerError {
			t.Errorf("Expected 200 or 500, got %d", rec.Code)
		}
		
		// Should return meaningful message
		if rec.Body.String() == "" {
			t.Error("Expected message in response body")
		}
		
		// Test other endpoints also handle DB failure
		endpoints := []string{
			"/dtako/events/import",
			"/dtako/ferry/import",
		}
		
		for _, endpoint := range endpoints {
			req = httptest.NewRequest("POST", endpoint, bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			rec = httptest.NewRecorder()
			
			r.ServeHTTP(rec, req)
			
			if rec.Code == http.StatusOK {
				// If 200, check for success:false
				var result models.ImportResult
				if err := json.Unmarshal(rec.Body.Bytes(), &result); err != nil {
					t.Errorf("Failed to parse response for %s: %v", endpoint, err)
				} else if result.Success {
					t.Errorf("Expected success:false for %s when production DB is unavailable", endpoint)
				}
			} else if rec.Code != http.StatusInternalServerError {
				t.Errorf("Expected 200 or 500 for %s, got %d", endpoint, rec.Code)
			}
		}
		
		// Restore original env vars
		if origHost != "" {
			os.Setenv("PROD_DB_HOST", origHost)
		}
		if origPort != "" {
			os.Setenv("PROD_DB_PORT", origPort)
		}
		
		t.Log("Production DB failure handled gracefully")
	})
	
	t.Run("Query local data when production unavailable", func(t *testing.T) {
		// Queries should still work if local DB is available
		req := httptest.NewRequest("GET", "/dtako/rows", nil)
		rec := httptest.NewRecorder()
		
		r.ServeHTTP(rec, req)
		
		if rec.Code != http.StatusOK {
			t.Errorf("Expected 200 for local query, got %d", rec.Code)
		}
		
		var rows []models.DtakoRow
		if err := json.Unmarshal(rec.Body.Bytes(), &rows); err != nil {
			t.Errorf("Failed to parse local query response: %v", err)
		}
		
		t.Logf("Local query returned %d rows", len(rows))
	})
}
