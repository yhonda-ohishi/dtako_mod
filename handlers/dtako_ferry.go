package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/yhonda-ohishi/dtako_mod/models"
	"github.com/yhonda-ohishi/dtako_mod/services"
)

// DtakoFerryHandler handles dtako_ferry related requests
type DtakoFerryHandler struct {
	service *services.DtakoFerryService
}

// NewDtakoFerryHandler creates a new dtako_ferry handler
func NewDtakoFerryHandler() *DtakoFerryHandler {
	return &DtakoFerryHandler{
		service: services.NewDtakoFerryService(),
	}
}

// List returns all dtako_ferry records
// @Summary List dtako ferry records
// @Description Get list of dtako ferry records with optional filters
// @Tags dtako_ferry
// @Accept json
// @Produce json
// @Param from query string false "From date (YYYY-MM-DD)"
// @Param to query string false "To date (YYYY-MM-DD)"
// @Param route query string false "Route filter"
// @Success 200 {array} models.DtakoFerry
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/dtako/ferry [get]
func (h *DtakoFerryHandler) List(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	route := r.URL.Query().Get("route")
	
	records, err := h.service.GetFerryRecords(from, to, route)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(records)
}

// Import imports dtako_ferry data from production
// @Summary Import dtako ferry data from production
// @Description Import dtako ferry data from production database
// @Tags dtako_ferry
// @Accept json
// @Produce json
// @Param request body models.ImportRequest true "Import request"
// @Success 200 {object} models.ImportResult
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/dtako/ferry/import [post]
func (h *DtakoFerryHandler) Import(w http.ResponseWriter, r *http.Request) {
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

	result, err := h.service.ImportFromProduction(req.FromDate, req.ToDate, req.Route)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// GetByID returns a specific dtako_ferry record by ID
// @Summary Get dtako ferry record by ID
// @Description Get a specific dtako ferry record by its ID
// @Tags dtako_ferry
// @Accept json
// @Produce json
// @Param id path string true "Ferry record ID"
// @Success 200 {object} models.DtakoFerry
// @Failure 404 {string} string "Not Found"
// @Router /api/dtako/ferry/{id} [get]
func (h *DtakoFerryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	record, err := h.service.GetFerryRecordByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(record)
}