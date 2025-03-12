package models

import "time"

// InstructorRequest is used for creating/updating an instructor.
type InstructorRequest struct {
	Name string `json:"name"`
}

// Instructor represents the instructor model.
type Instructor struct {
	InstructorID string    `json:"instructor_id"`
	UserID       string    `json:"user_id"`
	Name         string    `json:"name"`
	DateCreated  time.Time `json:"date_created"`
}
