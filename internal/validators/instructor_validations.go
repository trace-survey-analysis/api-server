package validators

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

// ValidateInstructorID checks if the provided string is a valid UUID
func ValidateInstructorID(instructorID string) error {
	if instructorID == "" {
		return errors.New("instructor ID is required")
	}

	// Check if the UUID is valid
	if _, err := uuid.Parse(instructorID); err != nil {
		return errors.New("invalid UUID format")
	}

	return nil
}

// ValidateInstructorName validates the instructor's name.
func ValidateInstructorName(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("name cannot be empty or just blank")
	}
	if ContainsNumber(name) {
		return errors.New("name should not contain any numbers")
	}
	return nil
}

// ValidateInstructorRequestBody validates if the instructor request body is valid
func ValidateInstructorRequestBody(contentLength int64, requireBody bool) error {
	if requireBody && contentLength == 0 {
		return errors.New("request body is required")
	}
	if !requireBody && contentLength > 0 {
		return errors.New("request body is not allowed")
	}
	return nil
}

// ValidateInstructorPatchFields validates the fields for patch requests
func ValidateInstructorPatchFields(fields map[string]interface{}) error {
	// Check if name field is provided
	nameVal, ok := fields["name"]
	if !ok {
		return errors.New("no valid fields to update")
	}

	// Check if name is a string
	name, ok := nameVal.(string)
	if !ok {
		return errors.New("invalid name format")
	}

	// Validate the name value
	return ValidateInstructorName(name)
}

// ExtractInstructorID extracts the instructor_id from the URL path.
// Expected path: /v1/instructor/{instructor_id}
func ExtractInstructorID(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) < 4 || parts[3] == "" {
		return ""
	}
	return parts[3]
}
