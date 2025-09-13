package integration

import (
	"github.com/go-chi/chi/v5"
	"github.com/yhonda-ohishi/dtako_mod"
)

// SetupTestRouter creates a test router with dtako routes mounted at /dtako
func SetupTestRouter() *chi.Mux {
	r := chi.NewRouter()

	// Mount dtako_mod routes at /dtako prefix for testing
	r.Route("/dtako", func(r chi.Router) {
		dtako_mod.RegisterRoutes(r)
	})

	return r
}