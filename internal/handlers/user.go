package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"strings"
	"unicode"

	"api-server/internal/database"
	"api-server/internal/models"
	"api-server/internal/repositories"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func extractUserID(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) < 4 || parts[3] == "" {
		return ""
	}
	return parts[3]
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	userID := extractUserID(r.URL.Path)
	if userID == "" {
		http.NotFound(w, r)
		return
	}

	// Check if the UUID is valid
	if _, err := uuid.Parse(userID); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid UUID format")
		return
	}

	switch r.Method {
	case http.MethodGet:
		GetUserHandler(w, r, userID)
	case http.MethodPut:
		UpdateUserHandler(w, r, userID)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func GetUserHandler(w http.ResponseWriter, r *http.Request, userID string) {
	// Query parameters check
	if len(r.URL.Query()) > 0 {
		respondWithError(w, http.StatusBadRequest, "query parameters are not allowed")
		return
	}

	// Body check
	if r.ContentLength > 0 {
		respondWithError(w, http.StatusBadRequest, "request body is not allowed")
		return
	}

	user, err := repositories.GetUserByID(database.GetDB(), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNoContent)
		} else {
			log.Printf("Error retrieving user: %v", err)
			w.WriteHeader(http.StatusServiceUnavailable)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if len(r.URL.Query()) > 0 {
		respondWithError(w, http.StatusBadRequest, "query parameters are not allowed")
		return
	}

	var req models.UserRequest

	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&req)
	}

	if err := validateUserRequest(req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	existingUser, err := repositories.GetUserByUsername(database.GetDB(), req.Username)
	if err != nil {
		log.Printf("Error checking existing user: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if existingUser != nil {
		respondWithError(w, http.StatusConflict, "user already exists")
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	req.Password = string(hashedPassword)

	user, err := repositories.CreateUser(database.GetDB(), req)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request, userID string) {
	// Check for query parameters
	if len(r.URL.Query()) > 0 {
		respondWithError(w, http.StatusBadRequest, "query parameters are not allowed")
		return
	}

	// Check if request empty
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && r.ContentLength > 0 {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Prevent username updates via body or parameters
	if _, ok := req["username"]; ok {
		respondWithError(w, http.StatusBadRequest, "username cannot be changed")
		return
	}

	user, err := repositories.GetUserByID(database.GetDB(), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "user not found")
		} else {
			log.Printf("Error retrieving user: %v", err)
			w.WriteHeader(http.StatusServiceUnavailable)
		}
		return
	}

	// Update only provided fields
	if firstName, ok := req["first_name"].(string); ok {
		if strings.TrimSpace(firstName) == "" || containsNumber(firstName) {
			respondWithError(w, http.StatusBadRequest, "invalid first_name")
			return
		}
		user.FirstName = firstName
	}
	if lastName, ok := req["last_name"].(string); ok {
		if strings.TrimSpace(lastName) == "" || containsNumber(lastName) {
			respondWithError(w, http.StatusBadRequest, "invalid last_name")
			return
		}
		user.LastName = lastName
	}

	if password, ok := req["password"].(string); ok {
		if len(password) < 8 {
			respondWithError(w, http.StatusBadRequest, "password must be at least 8 characters")
			return
		}

		// Hash the new password before updating it
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		user.Password = string(hashedPassword)
	}

	if err := repositories.UpdateUser(database.GetDB(), user); err != nil {
		log.Printf("Error updating user: %v", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func validateUserRequest(req models.UserRequest) error {
	if strings.TrimSpace(req.FirstName) == "" {
		return errors.New("first_name is required")
	}
	if strings.TrimSpace(req.LastName) == "" {
		return errors.New("last_name is required")
	}
	if strings.TrimSpace(req.Username) == "" {
		return errors.New("username is required")
	}
	if strings.TrimSpace(req.Password) == "" {
		return errors.New("password is required")
	}

	if containsNumber(req.FirstName) {
		return errors.New("first_name cannot contain numbers")
	}
	if containsNumber(req.LastName) {
		return errors.New("last_name cannot contain numbers")
	}

	if !isValidEmail(req.Username) {
		return errors.New("invalid email format")
	}

	if len(req.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	return nil
}

func containsNumber(s string) bool {
	for _, r := range s {
		if unicode.IsNumber(r) {
			return true
		}
	}
	return false
}

func isValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	return regexp.MustCompile(emailRegex).MatchString(email)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
