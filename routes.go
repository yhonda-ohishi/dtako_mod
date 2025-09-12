package dtako_mod

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/yhonda-ohishi/dtako_mod/handlers"
)

// RegisterRoutes registers all dtako_mod routes to the provided router
func RegisterRoutes(r chi.Router) {
	// Create handlers
	rowsHandler := handlers.NewDtakoRowsHandler()
	eventsHandler := handlers.NewDtakoEventsHandler()
	ferryHandler := handlers.NewDtakoFerryHandler()

	// Register routes
	r.Route("/dtako", func(r chi.Router) {
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

		// dtako_ferry endpoints
		r.Route("/ferry", func(r chi.Router) {
			r.Get("/", ferryHandler.List)
			r.Post("/import", ferryHandler.Import)
			r.Get("/{id}", ferryHandler.GetByID)
		})
	})
}

// Handler interface that each handler must implement
type Handler interface {
	List(w http.ResponseWriter, r *http.Request)
	Import(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
}