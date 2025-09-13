package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/yhonda-ohishi/dtako_mod/models"
	"github.com/yhonda-ohishi/dtako_mod/services"
)

// DtakoEventsHandler handles dtako_events related requests
type DtakoEventsHandler struct {
	service *services.DtakoEventsService
}

// NewDtakoEventsHandler creates a new dtako_events handler
func NewDtakoEventsHandler() *DtakoEventsHandler {
	return &DtakoEventsHandler{
		service: services.NewDtakoEventsService(),
	}
}

// List returns all dtako_events
// @Summary List dtako events
// @Description Get list of dtako events with optional filters
// @Tags dtako_events
// @Accept json
// @Produce json
// @Param from query string false "From date (YYYY-MM-DD)"
// @Param to query string false "To date (YYYY-MM-DD)"
// @Param type query string false "Event type filter"
// @Success 200 {array} models.DtakoEvent
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/dtako/events [get]
func (h *DtakoEventsHandler) List(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	eventType := r.URL.Query().Get("type")
	
	events, err := h.service.GetEvents(from, to, eventType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

// Import imports dtako_events from production
// @Summary Import dtako events from production
// @Description Import dtako events data from production database
// @Tags dtako_events
// @Accept json
// @Produce json
// @Param request body models.ImportRequest true "Import request"
// @Success 200 {object} models.ImportResult
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/dtako/events/import [post]
func (h *DtakoEventsHandler) Import(w http.ResponseWriter, r *http.Request) {
	var req models.ImportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set default date range if not provided
	if req.FromDate == "" {
		req.FromDate = time.Now().AddDate(0, -1, 0).Format("2006-01-02")
	}
	if req.ToDate == "" {
		req.ToDate = time.Now().Format("2006-01-02")
	}

	result, err := h.service.ImportFromProduction(req.FromDate, req.ToDate, req.EventType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// GetByID returns a specific dtako_event by ID
// @Summary Get dtako event by ID
// @Description Get a specific dtako event by its ID
// @Tags dtako_events
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {object} models.DtakoEvent
// @Failure 404 {string} string "Not Found"
// @Router /api/dtako/events/{id} [get]
func (h *DtakoEventsHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	event, err := h.service.GetEventByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}