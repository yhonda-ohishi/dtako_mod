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

// List lists dtako ferry records
// @Summary      List Dtako Ferry Records
// @Description  Get ferry operation data with optional filtering
// @Tags         dtako
// @Accept       json
// @Produce      json
// @Param        from    query     string  false  "Start date (YYYY-MM-DD)"
// @Param        to      query     string  false  "End date (YYYY-MM-DD)"
// @Param        route   query     string  false  "Route filter"
// @Success      200     {array}   models.DtakoFerry  "List of dtako ferry records"
// @Failure      400     {object}  models.ErrorResponse  "Invalid request parameters"
// @Failure      500     {object}  models.ErrorResponse  "Internal Server Error"
// @Router       /dtako/ferry [get]
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
// @Summary      Import Dtako Ferry Data
// @Description  Import ferry operation data from production database
// @Tags         dtako
// @Accept       json
// @Produce      json
// @Param        request body models.ImportRequest true "Import request"
// @Success      200     {object}  models.ImportResult  "Import successful"
// @Failure      400     {object}  models.ErrorResponse  "Bad Request"
// @Failure      500     {object}  models.ErrorResponse  "Internal Server Error"
// @Router       /dtako/ferry/import [post]
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
// @Summary      Get Dtako Ferry Record by ID
// @Description  Get specific ferry operation data by ID
// @Tags         dtako
// @Accept       json
// @Produce      json
// @Param        id      path      string  true  "Ferry record ID"
// @Success      200     {object}  models.DtakoFerry  "Dtako ferry record found"
// @Failure      404     {object}  models.ErrorResponse  "Not Found"
// @Router       /dtako/ferry/{id} [get]
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