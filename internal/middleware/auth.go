package middleware

import (
	"api-server/internal/database"
	"api-server/internal/repositories"
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// AuthMiddleware wraps handlers requiring authentication
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondWithError(w, http.StatusUnauthorized, "Authorization required")
			return
		}

		// Check if it's Basic auth
		if !strings.HasPrefix(authHeader, "Basic ") {
			respondWithError(w, http.StatusUnauthorized, "Invalid authorization method")
			return
		}

		// Decode credentials
		credentials, err := base64.StdEncoding.DecodeString(authHeader[6:])
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid authorization format")
			return
		}

		// Split username and password
		pair := strings.SplitN(string(credentials), ":", 2)
		if len(pair) != 2 {
			respondWithError(w, http.StatusUnauthorized, "Invalid authorization format")
			return
		}

		username := pair[0]
		password := pair[1]

		// Authenticate user
		user, err := authenticateUser(username, password)
		if err != nil || user == nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}

		// Check if the path is /v1/user/{user_id} and validate the user ID
		if strings.HasPrefix(r.URL.Path, "/v1/user/") {
			userIDFromPath := extractUserIDFromPath(r.URL.Path)
			if userIDFromPath == "" {
				respondWithError(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			if user.UserID != userIDFromPath {
				respondWithError(w, http.StatusForbidden, "Forbidden")
				return
			}
		}

		// Pass the authenticated user to the handler
		r = setUserContext(r, user)
		next(w, r)
	}
}

// authenticateUser validates username/password and returns user if valid
func authenticateUser(username, password string) (*repositories.UserWithPassword, error) {
	db := database.GetDB()

	// Get user with password
	user, err := repositories.GetUserWithPasswordByUsername(db, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Compare password with stored hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	return user, nil
}

// extractUserIDFromPath extracts user_id from URL path
func extractUserIDFromPath(path string) string {
	parts := strings.Split(path, "/")
	// For path like /v1/user/{user_id}
	if len(parts) < 4 || parts[3] == "" {
		return ""
	}
	return parts[3]
}

// Define a context key type to avoid collisions
type contextKey string

const userContextKey contextKey = "user"

// setUserContext adds user to request context
func setUserContext(r *http.Request, user *repositories.UserWithPassword) *http.Request {
	ctx := r.Context()
	return r.WithContext(context.WithValue(ctx, userContextKey, user))
}

// GetUserFromContext retrieves user from request context
func GetUserFromContext(r *http.Request) *repositories.UserWithPassword {
	if user, ok := r.Context().Value(userContextKey).(*repositories.UserWithPassword); ok {
		return user
	}
	return nil
}

// respondWithError sends a JSON error response with the provided status code
func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
