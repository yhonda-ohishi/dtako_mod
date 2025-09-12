package contract

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

// T010: Contract test POST /dtako/events/import
func TestPostDtakoEventsImport(t *testing.T) {
	// Setup router
	r := chi.NewRouter()
	dtako_mod.RegisterRoutes(r)

	// Test cases
	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		validateBody   func(*testing.T, []byte)
	}{
		{
			name: "Import events with date range",
			requestBody: models.ImportRequest{
				FromDate: "2025-01-01",
				ToDate:   "2025-01-31",
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body []byte) {
				var result models.ImportResult
				err := json.Unmarshal(body, &result)
				if err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				if result.Message == "" {
					t.Error("Expected message in response")
				}
				if result.ImportedAt.IsZero() {
					t.Error("Expected imported_at timestamp")
				}
			},
		},
		{
			name: "Import events with type filter",
			requestBody: map[string]interface{}{
				"from_date":   "2025-01-01",
				"to_date":     "2025-01-31",
				"event_type":  "ACCIDENT",
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body []byte) {
				var result models.ImportResult
				err := json.Unmarshal(body, &result)
				if err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				// Message should mention the event type
				if result.Message == "" {
					t.Error("Expected message with event type info")
				}
			},
		},
		{
			name: "Import with invalid event type",
			requestBody: map[string]interface{}{
				"from_date":   "2025-01-01",
				"to_date":     "2025-01-31",
				"event_type":  "INVALID_TYPE",
			},
			expectedStatus: http.StatusInternalServerError, // or BadRequest
			validateBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/dtako/events/import", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", 
					tt.expectedStatus, rec.Code, rec.Body.String())
			}

			if tt.validateBody != nil && rec.Code == http.StatusOK {
				tt.validateBody(t, rec.Body.Bytes())
			}
		})
	}
}