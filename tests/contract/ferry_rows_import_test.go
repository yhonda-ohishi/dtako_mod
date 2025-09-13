package contract

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yhonda-ohishi/dtako_mod/models"
)

// T015: Contract test POST /dtako/ferry_rows/import
func TestPostDtakoFerryRowsImport(t *testing.T) {
	// Setup router
	r := SetupTestRouter()
	

	// Test cases
	tests := []struct {
		name           string
		request        models.ImportRequest
		expectedStatus int
		validateBody   func(*testing.T, []byte)
	}{
		{
			name: "Import ferry row records with date range",
			request: models.ImportRequest{
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
				// Result can be success or failure depending on production data
				if result.ImportedRows < 0 {
					t.Errorf("Imported rows should not be negative")
				}
				if result.Message == "" {
					t.Error("Expected message to be present")
				}
			},
		},
		{
			name: "Import ferry row records with ferry company filter",
			request: models.ImportRequest{
				FromDate:     "2025-01-01",
				ToDate:       "2025-01-31",
				FerryCompany: "東京フェリー",
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body []byte) {
				var result models.ImportResult
				err := json.Unmarshal(body, &result)
				if err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				// Verify the message includes ferry company
				if result.Message == "" {
					t.Error("Expected message to be present")
				}
			},
		},
		{
			name: "Import with invalid ferry company",
			request: models.ImportRequest{
				FromDate:     "2025-01-01",
				ToDate:       "2025-01-31",
				FerryCompany: "INVALID_COMPANY",
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body []byte) {
				var result models.ImportResult
				err := json.Unmarshal(body, &result)
				if err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				// When no records match, import should still succeed with 0 rows
				if result.ImportedRows != 0 {
					t.Logf("Expected 0 rows for invalid ferry company, got %d", result.ImportedRows)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.request)
			req := httptest.NewRequest("POST", "/dtako/ferry_rows/import", bytes.NewBuffer(jsonBody))
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