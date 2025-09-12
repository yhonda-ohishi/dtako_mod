package contract

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/yhonda-ohishi/dtako_mod"
	"github.com/yhonda-ohishi/dtako_mod/models"
)

// T006: Contract test GET /dtako/rows
func TestGetDtakoRows(t *testing.T) {
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
			name:           "Get rows without date range",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body []byte) {
				var rows []models.DtakoRow
				err := json.Unmarshal(body, &rows)
				if err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				// Should return array (can be empty)
				if rows == nil {
					t.Error("Expected array response, got nil")
				}
			},
		},
		{
			name:           "Get rows with date range",
			queryParams:    "?from=2025-01-01&to=2025-01-31",
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body []byte) {
				var rows []models.DtakoRow
				err := json.Unmarshal(body, &rows)
				if err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				// Check all rows are within date range
				for _, row := range rows {
					if row.Date.Before(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)) ||
						row.Date.After(time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC)) {
						t.Errorf("Row date %v is outside range", row.Date)
					}
				}
			},
		},
		{
			name:           "Invalid date format",
			queryParams:    "?from=invalid&to=2025-01-31",
			expectedStatus: http.StatusInternalServerError, // or BadRequest depending on implementation
			validateBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req := httptest.NewRequest("GET", "/dtako/rows"+tt.queryParams, nil)
			rec := httptest.NewRecorder()

			// Execute request
			r.ServeHTTP(rec, req)

			// Check status code
			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			// Validate body if needed
			if tt.validateBody != nil && rec.Code == http.StatusOK {
				tt.validateBody(t, rec.Body.Bytes())
			}
		})
	}
}