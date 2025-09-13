package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/yhonda-ohishi/dtako_mod/models"
	"github.com/yhonda-ohishi/dtako_mod/repositories"
	"github.com/yhonda-ohishi/dtako_mod/services"
)

// DtakoFerryRowsHandler handles ferry row related requests
type DtakoFerryRowsHandler struct {
	service *services.DtakoFerryRowsService
}

// NewDtakoFerryRowsHandler creates a new ferry rows handler
func NewDtakoFerryRowsHandler() *DtakoFerryRowsHandler {
	// Initialize database connections
	repositories.InitDatabases()

	return &DtakoFerryRowsHandler{
		service: services.NewDtakoFerryRowsService(),
	}
}

// List handles GET /dtako/ferry_rows
// @Summary      List ferry row records
// @Description  Retrieve ferry row records with optional date range and ferry company filter
// @Tags         ferry_rows
// @Accept       json
// @Produce      json
// @Param        from          query     string  false  "Start date (YYYY-MM-DD)"
// @Param        to            query     string  false  "End date (YYYY-MM-DD)"
// @Param        ferry_company query     string  false  "Filter by ferry company name"
// @Success      200           {array}   models.DtakoFerryRow
// @Failure      400           {object}  models.ErrorResponse
// @Failure      500           {object}  models.ErrorResponse
// @Router       /dtako/ferry_rows [get]
func (h *DtakoFerryRowsHandler) List(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	ferryCompany := r.URL.Query().Get("ferry_company")

	records, err := h.service.GetFerryRows(from, to, ferryCompany)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(records)
}

// GetByID handles GET /dtako/ferry_rows/{id}
// @Summary      Get ferry row record by ID
// @Description  Retrieve a specific ferry row record by its ID
// @Tags         ferry_rows
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Ferry row record ID"
// @Success      200  {object}  models.DtakoFerryRow
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /dtako/ferry_rows/{id} [get]
func (h *DtakoFerryRowsHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	record, err := h.service.GetFerryRowByID(id)
	if err != nil {
		if err.Error() == "ferry row record not found: "+id {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(record)
}

// Import handles POST /dtako/ferry_rows/import
// @Summary      Import ferry row records from production
// @Description  Import ferry row records from production database for a date range
// @Tags         ferry_rows
// @Accept       json
// @Produce      json
// @Param        request  body      models.ImportRequest  true  "Import request with date range and optional ferry company filter"
// @Success      200      {object}  models.ImportResult
// @Failure      400      {object}  models.ErrorResponse
// @Failure      500      {object}  models.ErrorResponse
// @Router       /dtako/ferry_rows/import [post]
func (h *DtakoFerryRowsHandler) Import(w http.ResponseWriter, r *http.Request) {
	var req models.ImportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.FromDate == "" || req.ToDate == "" {
		http.Error(w, "from_date and to_date are required", http.StatusBadRequest)
		return
	}

	result, err := h.service.ImportFromProduction(req.FromDate, req.ToDate, req.FerryCompany)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}