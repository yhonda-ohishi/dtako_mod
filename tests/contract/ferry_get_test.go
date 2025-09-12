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

// T014: Contract test GET /dtako/ferry/{id}
func TestGetDtakoFerryByID(t *testing.T) {
	// Setup router
	r := chi.NewRouter()
	dtako_mod.RegisterRoutes(r)

	// Test cases
	tests := []struct {
		name           string
		id             string
		expectedStatus int
		validateBody   func(*testing.T, []byte)
	}{
		{
			name:           "Get existing ferry record",
			id:             "FERRY001",
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body []byte) {
				var record models.DtakoFerry
				err := json.Unmarshal(body, &record)
				if err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				if record.ID != "FERRY001" {
					t.Errorf("Expected ID FERRY001, got %s", record.ID)
				}
				// Validate required fields
				if record.Route == "" {
					t.Error("Expected route to be present")
				}
				if record.VehicleNo == "" {
					t.Error("Expected vehicle_no to be present")
				}
				if record.DriverCode == "" {
					t.Error("Expected driver_code to be present")
				}
				// Validate numeric fields
				if record.Passengers < 0 {
					t.Error("Passengers should not be negative")
				}
				if record.Vehicles < 0 {
					t.Error("Vehicles should not be negative")
				}
			},
		},
		{
			name:           "Get non-existent ferry record",
			id:             "NONEXISTENT",
			expectedStatus: http.StatusNotFound,
			validateBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/dtako/ferry/"+tt.id, nil)
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