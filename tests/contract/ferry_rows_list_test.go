package contract

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yhonda-ohishi/dtako_mod/models"
)

// T016: Contract test GET /dtako/ferry_rows
func TestGetDtakoFerryRows(t *testing.T) {
	// Setup router
	r := SetupTestRouter()
	

	// Test cases
	tests := []struct {
		name           string
		query          string
		expectedStatus int
		validateBody   func(*testing.T, []byte)
	}{
		{
			name:           "Get ferry row records without filters",
			query:          "",
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body []byte) {
				var records []models.DtakoFerryRow
				err := json.Unmarshal(body, &records)
				if err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
			},
		},
		{
			name:           "Get ferry row records with date range",
			query:          "?from=2024-01-15&to=2024-01-16",
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body []byte) {
				var records []models.DtakoFerryRow
				err := json.Unmarshal(body, &records)
				if err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				if len(records) == 0 {
					t.Logf("Warning: No records found for date range")
				}
			},
		},
		{
			name:           "Get ferry row records filtered by ferry company",
			query:          "?ferry_company=東京フェリー",
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body []byte) {
				var records []models.DtakoFerryRow
				err := json.Unmarshal(body, &records)
				if err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				// Verify all records have the specified ferry company
				for _, record := range records {
					if record.FerryCompanyName != "東京フェリー" {
						t.Errorf("Expected ferry company '東京フェリー', got '%s'", record.FerryCompanyName)
					}
				}
			},
		},
		{
			name:           "Get ferry row records with date range and ferry company filter",
			query:          "?from=2024-01-15&to=2024-01-16&ferry_company=東京フェリー",
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body []byte) {
				var records []models.DtakoFerryRow
				err := json.Unmarshal(body, &records)
				if err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				// Verify all records match the criteria
				for _, record := range records {
					if record.FerryCompanyName != "東京フェリー" {
						t.Errorf("Expected ferry company '東京フェリー', got '%s'", record.FerryCompanyName)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/dtako/ferry_rows"+tt.query, nil)
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