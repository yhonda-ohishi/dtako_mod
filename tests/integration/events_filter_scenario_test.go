package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yhonda-ohishi/dtako_mod/models"
)

// T016: Integration test - Filter dtako_events by type
func TestFilterEventsScenario(t *testing.T) {
	// Setup
	r := SetupTestRouter()

	t.Run("Filter events by type workflow", func(t *testing.T) {
		// Step 1: Import all types of events
		importReq := models.ImportRequest{
			FromDate: "2025-01-01",
			ToDate:   "2025-01-31",
		}
		
		bodyBytes, _ := json.Marshal(importReq)
		req := httptest.NewRequest("POST", "/dtako/events/import", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		
		r.ServeHTTP(rec, req)
		
		if rec.Code != http.StatusOK {
			t.Fatalf("Import failed with status %d: %s", rec.Code, rec.Body.String())
		}
		
		// Step 2: Query events without filter
		req = httptest.NewRequest("GET", "/dtako/events?from=2025-01-01&to=2025-01-31", nil)
		rec = httptest.NewRecorder()
		
		r.ServeHTTP(rec, req)
		
		if rec.Code != http.StatusOK {
			t.Fatalf("Query all events failed with status %d", rec.Code)
		}
		
		var allEvents []models.DtakoEvent
		json.Unmarshal(rec.Body.Bytes(), &allEvents)
		totalCount := len(allEvents)
		
		// Step 3: Query filtered by ACCIDENT type
		req = httptest.NewRequest("GET", "/dtako/events?from=2025-01-01&to=2025-01-31&type=ACCIDENT", nil)
		rec = httptest.NewRecorder()
		
		r.ServeHTTP(rec, req)
		
		if rec.Code != http.StatusOK {
			t.Fatalf("Query ACCIDENT events failed with status %d", rec.Code)
		}
		
		var accidentEvents []models.DtakoEvent
		json.Unmarshal(rec.Body.Bytes(), &accidentEvents)
		
		// Verify all returned events are ACCIDENT type
		for _, event := range accidentEvents {
			if event.EventType != "ACCIDENT" {
				t.Errorf("Expected ACCIDENT type, got %s", event.EventType)
			}
		}
		
		// Step 4: Query filtered by START type
		req = httptest.NewRequest("GET", "/dtako/events?from=2025-01-01&to=2025-01-31&type=START", nil)
		rec = httptest.NewRecorder()
		
		r.ServeHTTP(rec, req)
		
		if rec.Code != http.StatusOK {
			t.Fatalf("Query START events failed with status %d", rec.Code)
		}
		
		var startEvents []models.DtakoEvent
		json.Unmarshal(rec.Body.Bytes(), &startEvents)
		
		// Verify all returned events are START type
		for _, event := range startEvents {
			if event.EventType != "START" {
				t.Errorf("Expected START type, got %s", event.EventType)
			}
		}
		
		// Step 5: Import events with specific type filter
		importWithFilter := map[string]interface{}{
			"from_date":  "2025-02-01",
			"to_date":    "2025-02-28",
			"event_type": "運転",
		}
		
		bodyBytes, _ = json.Marshal(importWithFilter)
		req = httptest.NewRequest("POST", "/dtako/events/import", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rec = httptest.NewRecorder()
		
		r.ServeHTTP(rec, req)
		
		if rec.Code != http.StatusOK {
			t.Fatalf("Import with filter failed with status %d", rec.Code)
		}
		
		// Verify imported events
		req = httptest.NewRequest("GET", "/dtako/events?from=2025-02-01&to=2025-02-28&type=運転", nil)
		rec = httptest.NewRecorder()
		
		r.ServeHTTP(rec, req)
		
		var drivingEvents []models.DtakoEvent
		json.Unmarshal(rec.Body.Bytes(), &drivingEvents)

		for _, event := range drivingEvents {
			if event.EventType != "運転" {
				t.Errorf("Expected 運転 type, got %s", event.EventType)
			}
		}
		
		t.Logf("Total events: %d, Accident: %d, Start: %d, Driving: %d",
			totalCount, len(accidentEvents), len(startEvents), len(drivingEvents))
	})
}