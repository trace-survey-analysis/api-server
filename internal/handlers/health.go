package handlers

import (
	"api-server/internal/database"
	"api-server/internal/repositories"
	"log"
	"net/http"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET method
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Check for query parameters in the URL
	if len(r.URL.Query()) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request: Parameters are not allowed"))
		return
	}

	// Check for payload (body content)
	if r.ContentLength > 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request: Body is not allowed"))
		return
	}

	// Set headers
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	// Attempt to connect to DB
	db := database.GetDB()
	if db == nil {
		log.Printf("Health check failed: Database connection unavailable")
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("Service Unavailable"))
		return
	}

	// Use the repository to insert a health check record
	if err := repositories.InsertHealthCheck(db); err != nil {
		log.Printf("Health check failed: %v", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("Service Unavailable"))
		return
	}

	// Return 200 OK if everything works fine
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
