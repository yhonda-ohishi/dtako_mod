package contract

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yhonda-ohishi/dtako_mod/models"
)

// T007: Contract test POST /dtako/rows/import
func TestPostDtakoRowsImport(t *testing.T) {
	// Setup router
	r := SetupTestRouter()
	

	// Test cases
	tests := []struct {
		name           string
		requestBody    models.ImportRequest
		expectedStatus int
		validateBody   func(*testing.T, []byte)
	}{
		{
			name: "Import with valid date range",
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
				// Check required fields
				if result.Message == "" {
					t.Error("Expected message in response")
				}
				if result.ImportedAt.IsZero() {
					t.Error("Expected imported_at timestamp")
				}
			},
		},
		{
			name: "Import without date range (uses defaults)",
			requestBody: models.ImportRequest{
				FromDate: "",
				ToDate:   "",
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body []byte) {
				var result models.ImportResult
				err := json.Unmarshal(body, &result)
				if err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				// Should use default date range
				if result.Message == "" {
					t.Error("Expected message with default date range")
				}
			},
		},
		{
			name: "Import with invalid date format",
			requestBody: models.ImportRequest{
				FromDate: "invalid-date",
				ToDate:   "2025-01-31",
			},
			expectedStatus: http.StatusInternalServerError,
			validateBody:   nil,
		},
		{
			name: "Import with from_date after to_date",
			requestBody: models.ImportRequest{
				FromDate: "2025-02-01",
				ToDate:   "2025-01-01",
			},
			expectedStatus: http.StatusInternalServerError,
			validateBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare request body
			bodyBytes, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/dtako/rows/import", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			// Execute request
			r.ServeHTTP(rec, req)

			// Check status code
			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", 
					tt.expectedStatus, rec.Code, rec.Body.String())
			}

			// Validate body if needed
			if tt.validateBody != nil && rec.Code == http.StatusOK {
				tt.validateBody(t, rec.Body.Bytes())
			}
		})
	}
}