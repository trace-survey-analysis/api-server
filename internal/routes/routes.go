package routes

import (
	"github.com/gorilla/mux"

	"api-server/internal/handlers"
)

// Register all the application routes
func RegisterRoutes() *mux.Router {
	r := mux.NewRouter()

	// Public routes
	r.HandleFunc("/healthz", handlers.HealthCheckHandler).Methods("GET")

	return r
}
