package contract

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yhonda-ohishi/dtako_mod/models"
)

// T009: Contract test GET /dtako/events
func TestGetDtakoEvents(t *testing.T) {
	// Setup router
	r := SetupTestRouter()
	

	// Test cases
	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		validateBody   func(*testing.T, []byte)
	}{
		{
			name:           "Get events without filters",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body []byte) {
				var events []models.DtakoEvent
				err := json.Unmarshal(body, &events)
				if err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				if events == nil {
					t.Error("Expected array response, got nil")
				}
			},
		},
		{
			name:           "Get events with date range",
			queryParams:    "?from=2025-01-01&to=2025-01-31",
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body []byte) {
				var events []models.DtakoEvent
				err := json.Unmarshal(body, &events)
				if err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
			},
		},
		{
			name:           "Get events filtered by type",
			queryParams:    "?type=ACCIDENT",
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body []byte) {
				var events []models.DtakoEvent
				err := json.Unmarshal(body, &events)
				if err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				// All events should be of type ACCIDENT
				for _, event := range events {
					if event.EventType != "ACCIDENT" {
						t.Errorf("Expected event type ACCIDENT, got %s", event.EventType)
					}
				}
			},
		},
		{
			name:           "Get events with date range and type filter",
			queryParams:    "?from=2025-01-01&to=2025-01-31&type=START",
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body []byte) {
				var events []models.DtakoEvent
				err := json.Unmarshal(body, &events)
				if err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				for _, event := range events {
					if event.EventType != "START" {
						t.Errorf("Expected event type START, got %s", event.EventType)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/dtako/events"+tt.queryParams, nil)
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if tt.validateBody != nil && rec.Code == http.StatusOK {
				tt.validateBody(t, rec.Body.Bytes())
			}
		})
	}
}