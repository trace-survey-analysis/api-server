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
	r.HandleFunc("/v1/course/{course_id}", handlers.GetCourseHandler).Methods("GET")

	// Private routes
	//user
	r.HandleFunc("/v1/user/{user_id}", middleware.AuthMiddleware(handlers.UserHandler)).Methods("GET", "PUT")
	//instructor
	r.HandleFunc("/v1/instructor", middleware.AuthMiddleware(handlers.CreateInstructorHandler)).Methods("POST")
	r.HandleFunc("/v1/instructor/{instructor_id}", middleware.AuthMiddleware(handlers.InstructorHandler)).Methods("PUT", "PATCH", "DELETE")
	r.HandleFunc("/v1/instructors", middleware.AuthMiddleware(handlers.GetAllInstructorsHandler)).Methods("GET")
	//course
	r.HandleFunc("/v1/course", middleware.AuthMiddleware(handlers.CreateCourseHandler)).Methods("POST")
	r.HandleFunc("/v1/course/{course_id}", middleware.AuthMiddleware(handlers.CourseHandler)).Methods("PUT", "PATCH", "DELETE")
	r.HandleFunc("/v1/courses", middleware.AuthMiddleware(handlers.GetAllCoursesHandler)).Methods("GET")
	//trace
	r.HandleFunc("/v1/course/{course_id}/trace", middleware.AuthMiddleware(handlers.TraceHandler)).Methods("POST", "GET")
	r.HandleFunc("/v1/course/{course_id}/trace/{trace_id}", middleware.AuthMiddleware(handlers.TraceEntityHandler)).Methods("GET", "DELETE")
	r.HandleFunc("/v1/traces", middleware.AuthMiddleware(handlers.GetAllTracesHandler)).Methods("GET")

	return r
}
