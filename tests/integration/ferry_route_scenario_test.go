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

// T017: Integration test - Filter dtako_ferry by route
func TestFilterFerryByRouteScenario(t *testing.T) {
	// Setup
	r := chi.NewRouter()
	dtako_mod.RegisterRoutes(r)

	t.Run("Filter ferry records by route workflow", func(t *testing.T) {
		// Step 1: Import ferry data for all routes
		importReq := models.ImportRequest{
			FromDate: "2025-01-01",
			ToDate:   "2025-01-31",
		}
		
		bodyBytes, _ := json.Marshal(importReq)
		req := httptest.NewRequest("POST", "/dtako/ferry/import", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		
		r.ServeHTTP(rec, req)
		
		if rec.Code != http.StatusOK {
			t.Fatalf("Import failed with status %d: %s", rec.Code, rec.Body.String())
		}
		
		// Step 2: Query all ferry records
		req = httptest.NewRequest("GET", "/dtako/ferry?from=2025-01-01&to=2025-01-31", nil)
		rec = httptest.NewRecorder()
		
		r.ServeHTTP(rec, req)
		
		if rec.Code != http.StatusOK {
			t.Fatalf("Query all ferry records failed with status %d", rec.Code)
		}
		
		var allRecords []models.DtakoFerry
		json.Unmarshal(rec.Body.Bytes(), &allRecords)
		totalCount := len(allRecords)
		
		// Step 3: Query filtered by ROUTE_A
		req = httptest.NewRequest("GET", "/dtako/ferry?from=2025-01-01&to=2025-01-31&route=ROUTE_A", nil)
		rec = httptest.NewRecorder()
		
		r.ServeHTTP(rec, req)
		
		if rec.Code != http.StatusOK {
			t.Fatalf("Query ROUTE_A failed with status %d", rec.Code)
		}
		
		var routeARecords []models.DtakoFerry
		json.Unmarshal(rec.Body.Bytes(), &routeARecords)
		
		// Verify all returned records are for ROUTE_A
		for _, record := range routeARecords {
			if record.Route != "ROUTE_A" {
				t.Errorf("Expected ROUTE_A, got %s", record.Route)
			}
			// Verify arrival time is after departure time
			if !record.ArrivalTime.After(record.DepartureTime) {
				t.Error("Arrival time should be after departure time")
			}
			// Verify passenger and vehicle counts are non-negative
			if record.Passengers < 0 {
				t.Error("Passenger count should not be negative")
			}
			if record.Vehicles < 0 {
				t.Error("Vehicle count should not be negative")
			}
		}
		
		// Step 4: Query filtered by ROUTE_B
		req = httptest.NewRequest("GET", "/dtako/ferry?from=2025-01-01&to=2025-01-31&route=ROUTE_B", nil)
		rec = httptest.NewRecorder()
		
		r.ServeHTTP(rec, req)
		
		var routeBRecords []models.DtakoFerry
		json.Unmarshal(rec.Body.Bytes(), &routeBRecords)
		
		for _, record := range routeBRecords {
			if record.Route != "ROUTE_B" {
				t.Errorf("Expected ROUTE_B, got %s", record.Route)
			}
		}
		
		// Step 5: Import with specific route filter
		importWithRoute := map[string]interface{}{
			"from_date": "2025-02-01",
			"to_date":   "2025-02-28",
			"route":     "ROUTE_C",
		}
		
		bodyBytes, _ = json.Marshal(importWithRoute)
		req = httptest.NewRequest("POST", "/dtako/ferry/import", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rec = httptest.NewRecorder()
		
		r.ServeHTTP(rec, req)
		
		if rec.Code != http.StatusOK {
			t.Fatalf("Import with route filter failed with status %d", rec.Code)
		}
		
		// Verify imported records
		req = httptest.NewRequest("GET", "/dtako/ferry?from=2025-02-01&to=2025-02-28&route=ROUTE_C", nil)
		rec = httptest.NewRecorder()
		
		r.ServeHTTP(rec, req)
		
		var routeCRecords []models.DtakoFerry
		json.Unmarshal(rec.Body.Bytes(), &routeCRecords)
		
		for _, record := range routeCRecords {
			if record.Route != "ROUTE_C" {
				t.Errorf("Expected ROUTE_C, got %s", record.Route)
			}
		}
		
		t.Logf("Total records: %d, Route A: %d, Route B: %d, Route C: %d",
			totalCount, len(routeARecords), len(routeBRecords), len(routeCRecords))
	})
}