package validators

import (
	"errors"
	"strings"

	"api-server/internal/models"

	"github.com/google/uuid"
)

// validate Course code
func ValidateCourseCode(code string) error {
	if strings.TrimSpace(code) == "" {
		return errors.New("Course code cannot be empty or just blank")
	}
	if len(code) > 15 {
		return errors.New("Course code cannot be longer than 15 characters")
	}
	return nil
}

// validate Course name
func ValidateCourseName(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("Course name cannot be empty or just blank")
	}
	if len(name) > 100 {
		return errors.New("Course name cannot be longer than 100 characters")
	}
	return nil
}

// validate Course description
func ValidateCourseDescription(description string) error {
	if len(description) > 500 {
		return errors.New("Course description cannot be longer than 500 characters")
	}
	return nil
}

// validate Course instructor ID
func ValidateCourseInstructorID(instructorID string) error {
	if strings.TrimSpace(instructorID) == "" {
		return errors.New("Instructor ID cannot be empty or just blank")
	}
	if _, err := uuid.Parse(instructorID); err != nil {
		return errors.New("Invalid UUID format for instructor ID")
	}
	return nil
}

// validate Course department
func ValidateCourseDepartment(departmentID int) error {
	//departmentID cannot be empty
	if departmentID == 0 {
		return errors.New("Department ID cannot be empty")
	}
	return nil
}

// validate Course credit hours
func ValidateCourseCreditHours(creditHours int) error {
	if creditHours < 0 {
		return errors.New("Credit hours cannot be negative")
	}
	// credit hours cannot be higher than 4
	if creditHours > 4 {
		return errors.New("Credit hours cannot be higher than 4")
	}
	return nil
}

// validate Course request
func ValidateCourseRequest(courseReq models.CourseRequest) error {
	if err := ValidateCourseCode(courseReq.Code); err != nil {
		return err
	}
	if err := ValidateCourseName(courseReq.Name); err != nil {
		return err
	}
	if err := ValidateCourseDescription(courseReq.Description); err != nil {
		return err
	}
	if err := ValidateCourseInstructorID(courseReq.InstructorID); err != nil {
		return err
	}
	if err := ValidateCourseDepartment(courseReq.DepartmentID); err != nil {
		return err
	}

	if err := ValidateCourseCreditHours(courseReq.CreditHours); err != nil {
		return err
	}

	return nil
}

func ValidateCourseUpdateRequest(courseReq models.CourseRequest) error {
	if courseReq.Code != "" {
		if err := ValidateCourseCode(courseReq.Code); err != nil {
			return err
		}
	}
	if courseReq.Name != "" {
		if err := ValidateCourseName(courseReq.Name); err != nil {
			return err
		}
	}
	if courseReq.Description != "" {
		if err := ValidateCourseDescription(courseReq.Description); err != nil {
			return err
		}
	}
	if courseReq.CreditHours != 0 {
		if err := ValidateCourseCreditHours(courseReq.CreditHours); err != nil {
			return err
		}
	}
	if courseReq.InstructorID != "" {
		if err := ValidateCourseInstructorID(courseReq.InstructorID); err != nil {
			return err
		}
	}

	return nil
}
