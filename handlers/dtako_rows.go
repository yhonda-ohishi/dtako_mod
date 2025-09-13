package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/yhonda-ohishi/dtako_mod/models"
	"github.com/yhonda-ohishi/dtako_mod/services"
)

// DtakoRowsHandler handles dtako_rows related requests
type DtakoRowsHandler struct {
	service *services.DtakoRowsService
}

// NewDtakoRowsHandler creates a new dtako_rows handler
func NewDtakoRowsHandler() *DtakoRowsHandler {
	return &DtakoRowsHandler{
		service: services.NewDtakoRowsService(),
	}
}

// List returns all dtako_rows
// @Summary List dtako rows
// @Description Get list of dtako rows with optional date filter
// @Tags dtako_rows
// @Accept json
// @Produce json
// @Param from query string false "From date (YYYY-MM-DD)"
// @Param to query string false "To date (YYYY-MM-DD)"
// @Success 200 {array} models.DtakoRow
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/dtako/rows [get]
func (h *DtakoRowsHandler) List(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	
	rows, err := h.service.GetRows(from, to)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rows)
}

// Import imports dtako_rows from production
// @Summary Import dtako rows from production
// @Description Import dtako rows data from production database
// @Tags dtako_rows
// @Accept json
// @Produce json
// @Param request body models.ImportRequest true "Import request"
// @Success 200 {object} models.ImportResult
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/dtako/rows/import [post]
func (h *DtakoRowsHandler) Import(w http.ResponseWriter, r *http.Request) {
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

	result, err := h.service.ImportFromProduction(req.FromDate, req.ToDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// GetByID returns a specific dtako_row by ID
// @Summary Get dtako row by ID
// @Description Get a specific dtako row by its ID
// @Tags dtako_rows
// @Accept json
// @Produce json
// @Param id path string true "Row ID"
// @Success 200 {object} models.DtakoRow
// @Failure 404 {string} string "Not Found"
// @Router /api/dtako/rows/{id} [get]
func (h *DtakoRowsHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	row, err := h.service.GetRowByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(row)
}