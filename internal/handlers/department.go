package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"api-server/internal/database"
	"api-server/internal/repositories"
	"api-server/internal/validators"
)

// GetAllDepartmentsHandler handles GET /v1/departments
func GetAllDepartmentsHandler(w http.ResponseWriter, r *http.Request) {
	// Validate query parameters
	if err := validators.ValidateRequestParameters(r.URL.Query()); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate request body
	if err := validators.ValidateRequestBody(r.ContentLength, false); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	db := database.GetDB()
	departments, err := repositories.GetAllDepartments(db)
	if err != nil {
		log.Printf("Error retrieving departments: %v", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(departments)
}
