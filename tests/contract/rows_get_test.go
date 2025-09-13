package contract

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yhonda-ohishi/dtako_mod/models"
)

// T008: Contract test GET /dtako/rows/{id}
func TestGetDtakoRowByID(t *testing.T) {
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
			name:           "Get existing row",
			id:             "ROW001",
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body []byte) {
				var row models.DtakoRow
				err := json.Unmarshal(body, &row)
				if err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				if row.ID != "ROW001" {
					t.Errorf("Expected ID ROW001, got %s", row.ID)
				}
				// Validate required fields
				if row.VehicleNo == "" {
					t.Error("Expected vehicle_no to be present")
				}
				if row.DriverCode == "" {
					t.Error("Expected driver_code to be present")
				}
			},
		},
		{
			name:           "Get non-existent row",
			id:             "NONEXISTENT",
			expectedStatus: http.StatusNotFound,
			validateBody:   nil,
		},
		{
			name:           "Get with empty ID",
			id:             "",
			expectedStatus: http.StatusOK, // Empty ID matches list endpoint
			validateBody:   nil,
		},
	}

	// Insert test data (if needed for "existing row" test)
	// This would normally be done in test setup

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req := httptest.NewRequest("GET", "/dtako/rows/"+tt.id, nil)
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