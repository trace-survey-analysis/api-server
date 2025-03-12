package routes

import (
	"github.com/gorilla/mux"

	"api-server/internal/handlers"
	"api-server/internal/middleware"
)

// Register all the application routes
func RegisterRoutes() *mux.Router {
	r := mux.NewRouter()

	// Public routes
	r.HandleFunc("/healthz", handlers.HealthCheckHandler).Methods("GET")
	r.HandleFunc("/v1/user", handlers.CreateUserHandler).Methods("POST")
	r.HandleFunc("/v1/instructor/{instructor_id}", handlers.InstructorHandler).Methods("GET")

	// Private routes
	//user
	r.HandleFunc("/v1/user/{user_id}", middleware.AuthMiddleware(handlers.UserHandler)).Methods("GET", "PUT")
	//instructor
	r.HandleFunc("/v1/instructor", middleware.AuthMiddleware(handlers.CreateInstructorHandler)).Methods("POST")
	r.HandleFunc("/v1/instructor/{instructor_id}", middleware.AuthMiddleware(handlers.InstructorHandler)).Methods("PUT", "PATCH", "DELETE")

	return r
}
