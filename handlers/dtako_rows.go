// Package handlers provides API handlers for DTako module
//
//	@title			DTako API
//	@version		1.0.0
//	@description	Digital tachograph data management API for vehicle operation records
//	@BasePath		/dtako
//	@host			localhost:8080
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

// List lists dtako rows
// @Summary      List Dtako Rows
// @Description  Get vehicle operation data with optional date filtering
// @Tags         dtako_rows
// @Accept       json
// @Produce      json
// @Param        from    query     string  false  "Start date (YYYY-MM-DD)"
// @Param        to      query     string  false  "End date (YYYY-MM-DD)"
// @Success      200     {array}   models.DtakoRow  "List of dtako rows"
// @Failure      400     {object}  models.ErrorResponse  "Invalid request parameters"
// @Failure      500     {object}  models.ErrorResponse  "Internal Server Error"
// @Router       /rows [get]
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
// @Summary      Import Dtako Rows
// @Description  Import vehicle operation data from production database
// @Tags         dtako_rows
// @Accept       json
// @Produce      json
// @Param        request body models.ImportRequest true "Import request"
// @Success      200     {object}  models.ImportResult  "Import successful"
// @Failure      400     {object}  models.ErrorResponse  "Bad Request"
// @Failure      500     {object}  models.ErrorResponse  "Internal Server Error"
// @Router       /rows/import [post]
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
// @Summary      Get Dtako Row by ID
// @Description  Get specific vehicle operation data by ID
// @Tags         dtako_rows
// @Accept       json
// @Produce      json
// @Param        id      path      string  true  "Row ID"
// @Success      200     {object}  models.DtakoRow  "Dtako row found"
// @Failure      404     {object}  models.ErrorResponse  "Not Found"
// @Router       /rows/{id} [get]
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