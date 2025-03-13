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
	"api-server/internal/validators"

	"github.com/google/uuid"
)

func extractCourseID(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) < 4 || parts[3] == "" {
		return ""
	}
	return parts[3]
}

func CourseHandler(w http.ResponseWriter, r *http.Request) {
	courseID := extractCourseID(r.URL.Path)
	if courseID == "" {
		http.NotFound(w, r)
		return
	}

	if _, err := uuid.Parse(courseID); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid UUID format")
		return
	}

	switch r.Method {
	case http.MethodPut:
		UpdateCourseHandler(w, r, courseID)
	case http.MethodPatch:
		PatchCourseHandler(w, r, courseID)
	case http.MethodDelete:
		DeleteCourseHandler(w, r, courseID)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func CreateCourseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if len(r.URL.Query()) > 0 {
		respondWithError(w, http.StatusBadRequest, "query parameters are not allowed")
		return
	}

	var courseReq models.CourseRequest
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&courseReq)
	}
	// add course request validation from validators
	if err := validators.ValidateCourseRequest(courseReq); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get the user from the context
	user := middleware.GetUserFromContext(r)
	if user == nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID := user.UserID
	if _, err := repositories.GetInstructorByID(database.GetDB(), courseReq.InstructorID); err != nil {
		log.Printf("Error fetching instructor: %v", err)
		http.Error(w, "failed to get instructor", http.StatusBadRequest)
		return
	}

	if _, err := repositories.GetDepartmentByID(database.GetDB(), courseReq.DepartmentID); err != nil {
		log.Printf("Error fetching department: %v", err)
		http.Error(w, "failed to get department", http.StatusBadRequest)
		return
	}
	// TODO: add check for course code -> unique

	// Create the course model
	course := models.Course{
		CourseID:        uuid.New().String(),
		DateAdded:       time.Now().UTC(),
		DateLastUpdated: time.Now().UTC(),
		UserID:          userID,
		Code:            courseReq.Code,
		Name:            courseReq.Name,
		Description:     courseReq.Description,
		InstructorID:    courseReq.InstructorID,
		DepartmentID:    courseReq.DepartmentID,
		CreditHours:     courseReq.CreditHours,
	}

	log.Printf("Course: %v", course)
	// Create the course
	newCourse, err := repositories.CreateCourse(database.GetDB(), course)
	if err != nil {
		log.Printf("Error creating course: %v", err)
		http.Error(w, "failed to create course", http.StatusInternalServerError)
		return
	}
	// return 201 Created with course JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCourse)
}

// Get course by ID
func GetCourseHandler(w http.ResponseWriter, r *http.Request) {
	courseID := extractCourseID(r.URL.Path)
	if courseID == "" {
		http.Error(w, "invalid course ID", http.StatusBadRequest)
		return
	}
	course, err := repositories.GetCourseByID(database.GetDB(), courseID)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "course not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Error fetching course: %v", err)
		http.Error(w, "failed to get course", http.StatusInternalServerError)
		return
	}
	// return 200 OK with course JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(course)
}

// Update (PUT)	/v1/course/{course_id}
func UpdateCourseHandler(w http.ResponseWriter, r *http.Request, courseID string) {
	if courseID == "" {
		respondWithError(w, http.StatusBadRequest, "invalid course ID")
		return
	}

	// Do not allow query parameters.
	if len(r.URL.Query()) > 0 {
		respondWithError(w, http.StatusBadRequest, "query parameters are not allowed")
		return
	}

	// Decode request body.
	var req models.CourseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Getting existing values
	db := database.GetDB()
	existingCourse, err := repositories.GetCourseByID(db, courseID)
	if err != nil {
		log.Printf("Error fetching course: %v", err)
		respondWithError(w, http.StatusInternalServerError, "failed to fetch course")
		return
	}
	// update course request validation from validators
	if err := validators.ValidateCourseUpdateRequest(req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	// check if the instructor exists, if instructor_id is provided

	if req.InstructorID != "" {
		if _, err := repositories.GetInstructorByID(db, req.InstructorID); err != nil {
			log.Printf("Error fetching instructor: %v", err)
			respondWithError(w, http.StatusBadRequest, "instructor not found")
			return
		}
	}
	// check if the department exists
	if req.DepartmentID != 0 {
		if _, err := repositories.GetDepartmentByID(db, req.DepartmentID); err != nil {
			log.Printf("Error fetching department: %v", err)
			respondWithError(w, http.StatusBadRequest, "department not found")
			return
		}
	}
	// Create the course model
	course := models.Course{
		CourseID:        existingCourse.CourseID,
		DateAdded:       existingCourse.DateAdded,
		UserID:          existingCourse.UserID,
		DateLastUpdated: time.Now().UTC(),
		Code:            getValueOrDefault(req.Code, existingCourse.Code).(string), // Type assertion
		Name:            getValueOrDefault(req.Name, existingCourse.Name).(string),
		Description:     getValueOrDefault(req.Description, existingCourse.Description).(string),
		InstructorID:    getValueOrDefault(req.InstructorID, existingCourse.InstructorID).(string),
		DepartmentID:    getValueOrDefault(req.DepartmentID, existingCourse.DepartmentID).(int),
		CreditHours:     getValueOrDefault(req.CreditHours, existingCourse.CreditHours).(int), // Type assertion for int
	}

	if err := repositories.UpdateCourse(db, &course); err != nil {
		log.Printf("Error updating course: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(course)
}

// Patch (PATCH)	/v1/course/{course_id}
func PatchCourseHandler(w http.ResponseWriter, r *http.Request, courseID string) {
	if courseID == "" {
		respondWithError(w, http.StatusBadRequest, "invalid course ID")
		return
	}

	// Do not allow query parameters.
	if len(r.URL.Query()) > 0 {
		respondWithError(w, http.StatusBadRequest, "query parameters are not allowed")
		return
	}

	// Decode request body.
	var req models.CourseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Getting existing values
	db := database.GetDB()
	existingCourse, err := repositories.GetCourseByID(db, courseID)
	if err != nil {
		log.Printf("Error fetching course: %v", err)
		respondWithError(w, http.StatusInternalServerError, "failed to fetch course")
		return
	}
	// update course request validation from validators
	if err := validators.ValidateCourseUpdateRequest(req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	// check if the instructor exists, if instructor_id is provided

	if req.InstructorID != "" {
		if _, err := repositories.GetInstructorByID(db, req.InstructorID); err != nil {
			log.Printf("Error fetching instructor: %v", err)
			respondWithError(w, http.StatusBadRequest, "instructor not found")
			return
		}
	}
	// check if the department exists
	if req.DepartmentID != 0 {
		if _, err := repositories.GetDepartmentByID(db, req.DepartmentID); err != nil {
			log.Printf("Error fetching department: %v", err)
			respondWithError(w, http.StatusBadRequest, "department not found")
			return
		}
	}

	// Create the course model
	course := models.Course{
		CourseID:        existingCourse.CourseID,
		DateLastUpdated: time.Now().UTC(),
		DateAdded:       existingCourse.DateAdded,
		UserID:          existingCourse.UserID,
		Code:            getValueOrDefault(req.Code, existingCourse.Code).(string), // Type assertion
		Name:            getValueOrDefault(req.Name, existingCourse.Name).(string),
		Description:     getValueOrDefault(req.Description, existingCourse.Description).(string),
		InstructorID:    getValueOrDefault(req.InstructorID, existingCourse.InstructorID).(string),
		DepartmentID:    getValueOrDefault(req.DepartmentID, existingCourse.DepartmentID).(int),
		CreditHours:     getValueOrDefault(req.CreditHours, existingCourse.CreditHours).(int), // Type assertion for int
	}

	if err := repositories.UpdateCourse(db, &course); err != nil {
		log.Printf("Error updating course: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(course)
}

// Delete (DELETE)	/v1/course/{course_id}
func DeleteCourseHandler(w http.ResponseWriter, r *http.Request, courseID string) {
	if courseID == "" {
		respondWithError(w, http.StatusBadRequest, "invalid course ID")
		return
	}

	db := database.GetDB()
	err := repositories.DeleteCourse(db, courseID)
	if err != nil {
		log.Printf("Error deleting course: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Helper function to return the value from the request or the existing value if the request field is empty
func getValueOrDefault(newValue, defaultValue interface{}) interface{} {
	switch v := newValue.(type) {
	case string:
		if v != "" {
			return v
		}
	case int:
		if v > 0 {
			return v
		}
	}
	return defaultValue
}
