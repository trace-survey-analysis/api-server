package models

import "time"

// CourseRequest is the model for creating a new course
type CourseRequest struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	InstructorID string `json:"instructor_id"`
	DepartmentID int    `json:"department_id"`
	CreditHours  int    `json:"credit_hours"`
}
type Course struct {
	CourseID        string    `json:"course_id"`
	DateAdded       time.Time `json:"date_added"`
	DateLastUpdated time.Time `json:"date_last_updated"`
	UserID          string    `json:"user_id"`
	Code            string    `json:"code"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	InstructorID    string    `json:"instructor_id"`
	DepartmentID    int       `json:"department_id"`
	CreditHours     int       `json:"credit_hours"`
}
