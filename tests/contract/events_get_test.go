package contract

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yhonda-ohishi/dtako_mod/models"
)

// T011: Contract test GET /dtako/events/{id}
func TestGetDtakoEventByID(t *testing.T) {
	// Setup router with /dtako prefix
	r := SetupTestRouter()

	// Test cases
	tests := []struct {
		name           string
		id             string
		expectedStatus int
		validateBody   func(*testing.T, []byte)
	}{
		{
			name:           "Get existing event",
			id:             "EVENT001",
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body []byte) {
				var event models.DtakoEvent
				err := json.Unmarshal(body, &event)
				if err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				if event.ID != "EVENT001" {
					t.Errorf("Expected ID EVENT001, got %s", event.ID)
				}
				// Validate required fields
				if event.EventType == "" {
					t.Error("Expected event_type to be present")
				}
				if event.VehicleNo == "" {
					t.Error("Expected vehicle_no to be present")
				}
				if event.Description == "" {
					t.Error("Expected description to be present")
				}
			},
		},
		{
			name:           "Get non-existent event",
			id:             "NONEXISTENT",
			expectedStatus: http.StatusNotFound,
			validateBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/dtako/events/"+tt.id, nil)
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