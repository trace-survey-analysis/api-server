package validators

import (
	"api-server/internal/models"
	"errors"
	"regexp"
	"strings"
	"unicode"

	"github.com/google/uuid"
)

// ValidateUserID checks if the provided string is a valid UUID
func ValidateUserID(userID string) error {
	if userID == "" {
		return errors.New("user ID is required")
	}

	// Check if the UUID is valid
	if _, err := uuid.Parse(userID); err != nil {
		return errors.New("invalid UUID format")
	}

	return nil
}

// ValidateUserRequest validates all fields in the user creation request
func ValidateUserRequest(req models.UserRequest) error {
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

	if ContainsNumber(req.FirstName) {
		return errors.New("first_name cannot contain numbers")
	}
	if ContainsNumber(req.LastName) {
		return errors.New("last_name cannot contain numbers")
	}

	if !IsValidEmail(req.Username) {
		return errors.New("invalid email format")
	}

	if len(req.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	return nil
}

// ValidateUserUpdateFields validates the fields for user update
func ValidateUserUpdateFields(fields map[string]interface{}) error {
	// Prevent username updates
	if _, ok := fields["username"]; ok {
		return errors.New("username cannot be changed")
	}

	// Validate first_name if provided
	if firstName, ok := fields["first_name"].(string); ok {
		if strings.TrimSpace(firstName) == "" {
			return errors.New("first_name cannot be empty")
		}
		if ContainsNumber(firstName) {
			return errors.New("first_name cannot contain numbers")
		}
	}

	// Validate last_name if provided
	if lastName, ok := fields["last_name"].(string); ok {
		if strings.TrimSpace(lastName) == "" {
			return errors.New("last_name cannot be empty")
		}
		if ContainsNumber(lastName) {
			return errors.New("last_name cannot contain numbers")
		}
	}

	// Validate password if provided
	if password, ok := fields["password"].(string); ok {
		if len(password) < 8 {
			return errors.New("password must be at least 8 characters")
		}
	}

	return nil
}

// ValidateRequestParameters checks for unwanted query parameters
func ValidateRequestParameters(queryParams map[string][]string) error {
	if len(queryParams) > 0 {
		return errors.New("query parameters are not allowed")
	}
	return nil
}

// ValidateRequestBody checks if a request body is present when it shouldn't be
func ValidateRequestBody(contentLength int64, allowBody bool) error {
	if !allowBody && contentLength > 0 {
		return errors.New("request body is not allowed")
	}
	return nil
}

// ContainsNumber checks if a string contains any numeric characters
func ContainsNumber(s string) bool {
	for _, r := range s {
		if unicode.IsNumber(r) {
			return true
		}
	}
	return false
}

// IsValidEmail checks if a string is in valid email format
func IsValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	return regexp.MustCompile(emailRegex).MatchString(email)
}
