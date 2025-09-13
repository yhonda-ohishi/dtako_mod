package contract

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yhonda-ohishi/dtako_mod/models"
)

// T014: Contract test GET /dtako/ferry_rows/{id}
func TestGetDtakoFerryRowByID(t *testing.T) {
	// Setup router
	r := SetupTestRouter()
	

	// Test cases
	tests := []struct {
		name           string
		id             string
		expectedStatus int
		validateBody   func(*testing.T, []byte)
	}{
		{
			name:           "Get existing ferry row record",
			id:             "1",
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body []byte) {
				var record models.DtakoFerryRow
				err := json.Unmarshal(body, &record)
				if err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				if record.ID != 1 {
					t.Errorf("Expected ID 1, got %d", record.ID)
				}
				// Validate required fields
				if record.UnkoNo == "" {
					t.Error("Expected unko_no to be present")
				}
				if record.FerryCompanyName == "" {
					t.Error("Expected ferry_company_name to be present")
				}
				if record.BoardingName == "" {
					t.Error("Expected boarding_name to be present")
				}
				// Validate numeric fields
				if record.StandardFare < 0 {
					t.Error("Standard fare should not be negative")
				}
				if record.ContractFare < 0 {
					t.Error("Contract fare should not be negative")
				}
			},
		},
		{
			name:           "Get non-existent ferry row record",
			id:             "999999",
			expectedStatus: http.StatusNotFound,
			validateBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/dtako/ferry_rows/"+tt.id, nil)
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