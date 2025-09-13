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

// List lists dtako events
// @Summary      List Dtako Events
// @Description  Get event data with location information and optional filtering
// @Tags         dtako
// @Accept       json
// @Produce      json
// @Param        from     query     string  false  "Start date (YYYY-MM-DD)"
// @Param        to       query     string  false  "End date (YYYY-MM-DD)"
// @Param        type     query     string  false  "Event type filter"
// @Param        unko_no  query     string  false  "Filter by 運行NO (links to dtako_rows)"
// @Success      200      {array}   models.DtakoEvent  "List of dtako events"
// @Failure      400      {object}  models.ErrorResponse  "Invalid request parameters"
// @Failure      500      {object}  models.ErrorResponse  "Internal Server Error"
// @Router       /events [get]
func (h *DtakoEventsHandler) List(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	eventType := r.URL.Query().Get("type")
	unkoNo := r.URL.Query().Get("unko_no")

	events, err := h.service.GetEvents(from, to, eventType, unkoNo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

// Import imports dtako_events from production
// @Summary      Import Dtako Events
// @Description  Import event data from production database
// @Tags         dtako
// @Accept       json
// @Produce      json
// @Param        request body models.ImportRequest true "Import request"
// @Success      200     {object}  models.ImportResult  "Import successful"
// @Failure      400     {object}  models.ErrorResponse  "Bad Request"
// @Failure      500     {object}  models.ErrorResponse  "Internal Server Error"
// @Router       /events/import [post]
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
// @Summary      Get Dtako Event by ID
// @Description  Get specific event data by ID
// @Tags         dtako
// @Accept       json
// @Produce      json
// @Param        id      path      string  true  "Event ID"
// @Success      200     {object}  models.DtakoEvent  "Dtako event found"
// @Failure      404     {object}  models.ErrorResponse  "Not Found"
// @Router       /events/{id} [get]
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