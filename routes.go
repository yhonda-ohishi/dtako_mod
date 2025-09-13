package dtako_mod

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/yhonda-ohishi/dtako_mod/handlers"
)

// RegisterRoutes registers all dtako_mod routes to the provided router
// Note: Routes do not include /dtako prefix - this should be set by the parent router
func RegisterRoutes(r chi.Router) {
	// Initialize database connections
	// This is done automatically when handlers are created

	// Create handlers
	rowsHandler := handlers.NewDtakoRowsHandler()
	eventsHandler := handlers.NewDtakoEventsHandler()
	ferryRowsHandler := handlers.NewDtakoFerryRowsHandler()

	// Register routes WITHOUT /dtako prefix
	// dtako_rows endpoints
	r.Route("/rows", func(r chi.Router) {
		r.Get("/", rowsHandler.List)
		r.Post("/import", rowsHandler.Import)
		r.Get("/{id}", rowsHandler.GetByID)
	})

	// dtako_events endpoints
	r.Route("/events", func(r chi.Router) {
		r.Get("/", eventsHandler.List)
		r.Post("/import", eventsHandler.Import)
		r.Get("/{id}", eventsHandler.GetByID)
	})

	// dtako_ferry_rows endpoints
	r.Route("/ferry_rows", func(r chi.Router) {
		r.Get("/", ferryRowsHandler.List)
		r.Post("/import", ferryRowsHandler.Import)
		r.Get("/{id}", ferryRowsHandler.GetByID)
	})
}

// Handler interface that each handler must implement
type Handler interface {
	List(w http.ResponseWriter, r *http.Request)
	Import(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
}