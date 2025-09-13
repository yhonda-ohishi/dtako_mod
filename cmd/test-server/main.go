package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yhonda-ohishi/dtako_mod"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Mount dtako_mod routes under /dtako prefix
	r.Route("/dtako", func(r chi.Router) {
		dtako_mod.RegisterRoutes(r)
	})

	log.Println("Server starting on :8080")
	log.Println("Test endpoints:")
	log.Println("  GET http://localhost:8080/dtako/events?from=2025-09-13")
	log.Println("  GET http://localhost:8080/dtako/rows?from=2025-09-13")
	log.Println("  GET http://localhost:8080/dtako/ferry_rows?from=2025-09-13")

	http.ListenAndServe(":8080", r)
}