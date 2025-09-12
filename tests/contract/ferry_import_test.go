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

// T013: Contract test POST /dtako/ferry/import
func TestPostDtakoFerryImport(t *testing.T) {
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
			name: "Import ferry records with date range",
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
			name: "Import ferry records with route filter",
			requestBody: map[string]interface{}{
				"from_date": "2025-01-01",
				"to_date":   "2025-01-31",
				"route":     "ROUTE_A",
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body []byte) {
				var result models.ImportResult
				err := json.Unmarshal(body, &result)
				if err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				// Message should mention the route
				if result.Message == "" {
					t.Error("Expected message with route info")
				}
			},
		},
		{
			name: "Import with invalid route",
			requestBody: map[string]interface{}{
				"from_date": "2025-01-01",
				"to_date":   "2025-01-31",
				"route":     "INVALID_ROUTE",
			},
			expectedStatus: http.StatusInternalServerError, // or BadRequest
			validateBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/dtako/ferry/import", bytes.NewReader(bodyBytes))
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