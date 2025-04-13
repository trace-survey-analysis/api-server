package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"api-server/internal/database"
	"api-server/internal/middleware"
	"api-server/internal/models"
	"api-server/internal/repositories"
	"api-server/internal/validators"

	"github.com/google/uuid"
)

func InstructorHandler(w http.ResponseWriter, r *http.Request) {
	instructorID := validators.ExtractInstructorID(r.URL.Path)
	if instructorID == "" {
		http.NotFound(w, r)
		return
	}

	if err := validators.ValidateInstructorID(instructorID); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	switch r.Method {
	case http.MethodGet:
		GetInstructorHandler(w, r, instructorID)
	case http.MethodPut:
		UpdateInstructorHandler(w, r, instructorID)
	case http.MethodPatch:
		PatchInstructorHandler(w, r, instructorID)
	case http.MethodDelete:
		DeleteInstructorHandler(w, r, instructorID)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// CreateInstructorHandler handles POST /v1/instructor.
func CreateInstructorHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Validate query parameters
	if err := validators.ValidateRequestParameters(r.URL.Query()); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Decode request body.
	var req models.InstructorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate name.
	if err := validators.ValidateInstructorName(req.Name); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get authenticated user from context.
	authUser := middleware.GetUserFromContext(r)
	if authUser == nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Create the instructor model.
	instructor := models.Instructor{
		InstructorID: uuid.New().String(),
		UserID:       authUser.UserID,
		Name:         req.Name,
		DateCreated:  time.Now().UTC(),
	}

	db := database.GetDB()
	createdInstructor, err := repositories.CreateInstructor(db, instructor)
	if err != nil {
		log.Printf("Error creating instructor: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdInstructor)
}

// UpdateInstructorHandler handles PUT /v1/instructor/{instructor_id}.
func UpdateInstructorHandler(w http.ResponseWriter, r *http.Request, instructorID string) {
	// Validate query parameters
	if err := validators.ValidateRequestParameters(r.URL.Query()); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if instructorID == "" {
		http.NotFound(w, r)
		return
	}

	var req models.InstructorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || r.ContentLength == 0 {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validators.ValidateInstructorName(req.Name); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	db := database.GetDB()
	instructor, err := repositories.GetInstructorByID(db, instructorID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "instructor not found")
			return
		}
		log.Printf("Error retrieving instructor: %v", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// Update only the name.
	instructor.Name = req.Name
	if err := repositories.UpdateInstructor(db, instructor); err != nil {
		log.Printf("Error updating instructor: %v", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// PatchInstructorHandler handles PATCH /v1/instructor/{instructor_id}.
func PatchInstructorHandler(w http.ResponseWriter, r *http.Request, instructorID string) {
	// Validate query parameters
	if err := validators.ValidateRequestParameters(r.URL.Query()); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if instructorID == "" {
		http.NotFound(w, r)
		return
	}

	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || r.ContentLength == 0 {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate patch fields (name)
	if err := validators.ValidateInstructorPatchFields(req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	db := database.GetDB()
	instructor, err := repositories.GetInstructorByID(db, instructorID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "instructor not found")
			return
		}
		log.Printf("Error retrieving instructor: %v", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// Update the name from validated request
	instructor.Name = req["name"].(string)
	if err := repositories.UpdateInstructor(db, instructor); err != nil {
		log.Printf("Error patching instructor: %v", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DeleteInstructorHandler handles DELETE /v1/instructor/{instructor_id}.
func DeleteInstructorHandler(w http.ResponseWriter, r *http.Request, instructorID string) {
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

	if instructorID == "" {
		http.NotFound(w, r)
		return
	}

	db := database.GetDB()
	if err := repositories.DeleteInstructor(db, instructorID); err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "instructor not found")
			return
		}
		log.Printf("Error deleting instructor: %v", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetInstructorHandler handles GET /v1/instructor/{instructor_id}.
func GetInstructorHandler(w http.ResponseWriter, r *http.Request, instructorID string) {
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

	if instructorID == "" {
		http.NotFound(w, r)
		return
	}

	db := database.GetDB()
	instructor, err := repositories.GetInstructorByID(db, instructorID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "instructor not found")
			return
		}
		log.Printf("Error retrieving instructor: %v", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(instructor)
}

// GetAllInstructorsHandler handles GET /v1/instructors
func GetAllInstructorsHandler(w http.ResponseWriter, r *http.Request) {
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
	instructors, err := repositories.GetAllInstructors(db)
	if err != nil {
		log.Printf("Error retrieving instructors: %v", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(instructors)
}
