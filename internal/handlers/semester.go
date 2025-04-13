package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"api-server/internal/database"
	"api-server/internal/repositories"
	"api-server/internal/validators"
)

// GetAllSemesterTermsHandler handles GET /v1/semesters
func GetAllSemesterTermsHandler(w http.ResponseWriter, r *http.Request) {
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
	semesterTerms, err := repositories.GetAllSemesterTerms(db)
	if err != nil {
		log.Printf("Error retrieving semester terms: %v", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(semesterTerms)
}
