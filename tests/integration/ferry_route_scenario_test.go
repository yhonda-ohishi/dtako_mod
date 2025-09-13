package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yhonda-ohishi/dtako_mod/models"
)

// T017: Integration test - Filter dtako_ferry by route
func TestFilterFerryByRouteScenario(t *testing.T) {
	// Setup
	r := SetupTestRouter()

	t.Run("Filter ferry records by route workflow", func(t *testing.T) {
		// Step 1: Import ferry data for all routes
		importReq := models.ImportRequest{
			FromDate: "2025-01-01",
			ToDate:   "2025-01-31",
		}
		
		bodyBytes, _ := json.Marshal(importReq)
		req := httptest.NewRequest("POST", "/dtako/ferry_rows/import", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		
		r.ServeHTTP(rec, req)
		
		if rec.Code != http.StatusOK {
			t.Fatalf("Import failed with status %d: %s", rec.Code, rec.Body.String())
		}
		
		// Step 2: Query all ferry records
		req = httptest.NewRequest("GET", "/dtako/ferry_rows?from=2025-01-01&to=2025-01-31", nil)
		rec = httptest.NewRecorder()
		
		r.ServeHTTP(rec, req)
		
		if rec.Code != http.StatusOK {
			t.Fatalf("Query all ferry records failed with status %d", rec.Code)
		}
		
		var allRecords []models.DtakoFerryRow
		json.Unmarshal(rec.Body.Bytes(), &allRecords)
		totalCount := len(allRecords)
		
		// Step 3: Query filtered by ROUTE_A
		req = httptest.NewRequest("GET", "/dtako/ferry_rows?from=2025-01-01&to=2025-01-31&ferry_company=東京フェリー", nil)
		rec = httptest.NewRecorder()
		
		r.ServeHTTP(rec, req)
		
		if rec.Code != http.StatusOK {
			t.Fatalf("Query ROUTE_A failed with status %d", rec.Code)
		}
		
		var routeARecords []models.DtakoFerryRow
		json.Unmarshal(rec.Body.Bytes(), &routeARecords)
		
		// Verify all returned records are for 東京フェリー
		for _, record := range routeARecords {
			if record.FerryCompanyName != "東京フェリー" {
				t.Errorf("Expected 東京フェリー, got %s", record.FerryCompanyName)
			}
		}
		
		// Step 4: Query filtered by ROUTE_B
		req = httptest.NewRequest("GET", "/dtako/ferry_rows?from=2025-01-01&to=2025-01-31&ferry_company=大阪フェリー", nil)
		rec = httptest.NewRecorder()
		
		r.ServeHTTP(rec, req)
		
		var routeBRecords []models.DtakoFerryRow
		json.Unmarshal(rec.Body.Bytes(), &routeBRecords)
		
		for _, record := range routeBRecords {
			if record.FerryCompanyName != "大阪フェリー" {
				t.Errorf("Expected 大阪フェリー, got %s", record.FerryCompanyName)
			}
		}
		
		// Step 5: Import with specific ferry company filter
		importWithRoute := map[string]interface{}{
			"from_date": "2025-02-01",
			"to_date":   "2025-02-28",
			"ferry_company":     "神戸フェリー",
		}
		
		bodyBytes, _ = json.Marshal(importWithRoute)
		req = httptest.NewRequest("POST", "/dtako/ferry_rows/import", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rec = httptest.NewRecorder()
		
		r.ServeHTTP(rec, req)
		
		if rec.Code != http.StatusOK {
			t.Fatalf("Import with route filter failed with status %d", rec.Code)
		}
		
		// Verify imported records
		req = httptest.NewRequest("GET", "/dtako/ferry_rows?from=2025-02-01&to=2025-02-28&ferry_company=神戸フェリー", nil)
		rec = httptest.NewRecorder()
		
		r.ServeHTTP(rec, req)
		
		var routeCRecords []models.DtakoFerryRow
		json.Unmarshal(rec.Body.Bytes(), &routeCRecords)
		
		for _, record := range routeCRecords {
			if record.FerryCompanyName != "神戸フェリー" {
				t.Errorf("Expected 神戸フェリー, got %s", record.FerryCompanyName)
			}
		}
		
		t.Logf("Total records: %d, Route A: %d, Route B: %d, Route C: %d",
			totalCount, len(routeARecords), len(routeBRecords), len(routeCRecords))
	})
}