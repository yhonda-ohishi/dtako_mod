package contract

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/yhonda-ohishi/dtako_mod"
	"github.com/yhonda-ohishi/dtako_mod/models"
)

// T012: Contract test GET /dtako/ferry
func TestGetDtakoFerry(t *testing.T) {
	// Setup router
	r := chi.NewRouter()
	dtako_mod.RegisterRoutes(r)

	// Test cases
	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		validateBody   func(*testing.T, []byte)
	}{
		{
			name:           "Get ferry records without filters",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body []byte) {
				var records []models.DtakoFerry
				err := json.Unmarshal(body, &records)
				if err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				if records == nil {
					t.Error("Expected array response, got nil")
				}
			},
		},
		{
			name:           "Get ferry records with date range",
			queryParams:    "?from=2025-01-01&to=2025-01-31",
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body []byte) {
				var records []models.DtakoFerry
				err := json.Unmarshal(body, &records)
				if err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
			},
		},
		{
			name:           "Get ferry records filtered by route",
			queryParams:    "?route=ROUTE_A",
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body []byte) {
				var records []models.DtakoFerry
				err := json.Unmarshal(body, &records)
				if err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				// All records should be for ROUTE_A
				for _, record := range records {
					if record.Route != "ROUTE_A" {
						t.Errorf("Expected route ROUTE_A, got %s", record.Route)
					}
				}
			},
		},
		{
			name:           "Get ferry with date range and route filter",
			queryParams:    "?from=2025-01-01&to=2025-01-31&route=ROUTE_B",
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body []byte) {
				var records []models.DtakoFerry
				err := json.Unmarshal(body, &records)
				if err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				for _, record := range records {
					if record.Route != "ROUTE_B" {
						t.Errorf("Expected route ROUTE_B, got %s", record.Route)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/dtako/ferry"+tt.queryParams, nil)
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if tt.validateBody != nil && rec.Code == http.StatusOK {
				tt.validateBody(t, rec.Body.Bytes())
			}
		})
	}
}