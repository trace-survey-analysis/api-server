package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"api-server/internal/database"
	"api-server/internal/models"
	"api-server/internal/repositories"
	"api-server/internal/validators"

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

	// Validate the user ID
	if err := validators.ValidateUserID(userID); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
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

	// Validate query parameters
	if err := validators.ValidateRequestParameters(r.URL.Query()); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req models.UserRequest

	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&req)
	}

	// Validate user request
	if err := validators.ValidateUserRequest(req); err != nil {
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
	// Validate query parameters
	if err := validators.ValidateRequestParameters(r.URL.Query()); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Parse request body
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && r.ContentLength > 0 {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate update fields
	if err := validators.ValidateUserUpdateFields(req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
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
		user.FirstName = firstName
	}
	if lastName, ok := req["last_name"].(string); ok {
		user.LastName = lastName
	}

	if password, ok := req["password"].(string); ok {
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

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
