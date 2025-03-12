package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"api-server/internal/database"
	"api-server/internal/middleware"
	"api-server/internal/models"
	"api-server/internal/repositories"

	"github.com/google/uuid"
)

func InstructorHandler(w http.ResponseWriter, r *http.Request) {
	instructorID := extractInstructorID(r.URL.Path)
	if instructorID == "" {
		http.NotFound(w, r)
		return
	}

	if _, err := uuid.Parse(instructorID); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid UUID format")
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

	// Do not allow query parameters.
	if len(r.URL.Query()) > 0 {
		respondWithError(w, http.StatusBadRequest, "query parameters are not allowed")
		return
	}

	// Decode request body.
	var req models.InstructorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate name.
	if err := validateInstructorName(req.Name); err != nil {
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
	// Do not allow query parameters.
	if len(r.URL.Query()) > 0 {
		respondWithError(w, http.StatusBadRequest, "query parameters are not allowed")
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

	if err := validateInstructorName(req.Name); err != nil {
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
	// Do not allow query parameters.
	if len(r.URL.Query()) > 0 {
		respondWithError(w, http.StatusBadRequest, "query parameters are not allowed")
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

	// Only the "name" field is allowed.
	nameVal, ok := req["name"]
	if !ok {
		respondWithError(w, http.StatusBadRequest, "no valid fields to update")
		return
	}
	name, ok := nameVal.(string)
	if !ok {
		respondWithError(w, http.StatusBadRequest, "invalid name format")
		return
	}
	if err := validateInstructorName(name); err != nil {
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

	instructor.Name = name
	if err := repositories.UpdateInstructor(db, instructor); err != nil {
		log.Printf("Error patching instructor: %v", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DeleteInstructorHandler handles DELETE /v1/instructor/{instructor_id}.
func DeleteInstructorHandler(w http.ResponseWriter, r *http.Request, instructorID string) {
	// Do not allow query parameters.
	if len(r.URL.Query()) > 0 {
		respondWithError(w, http.StatusBadRequest, "query parameters are not allowed")
		return
	}
	// Do not allow a request body.
	if r.ContentLength > 0 {
		respondWithError(w, http.StatusBadRequest, "request body is not allowed")
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
	// Do not allow query parameters.
	if len(r.URL.Query()) > 0 {
		respondWithError(w, http.StatusBadRequest, "query parameters are not allowed")
		return
	}
	// Do not allow a request body.
	if r.ContentLength > 0 {
		respondWithError(w, http.StatusBadRequest, "request body is not allowed")
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

// extractInstructorID extracts the instructor_id from the URL path.
// Expected path: /v1/instructor/{instructor_id}
func extractInstructorID(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) < 4 || parts[3] == "" {
		return ""
	}
	return parts[3]
}

// validateInstructorName validates the instructor's name.
func validateInstructorName(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("name cannot be empty or just blank")
	}
	if containsNumber(name) {
		return errors.New("name should not contain any numbers")
	}
	return nil
}
