package models

import "time"

type TraceRequest struct {
	InstructorID string `json:"instructor_id"`
	SemesterTerm string `json:"semester_term"`
	Section      string `json:"section"`
}

type Trace struct {
	TraceID      string    `json:"trace_id"`
	UserID       string    `json:"user_id"`
	FileName     string    `json:"file_name"`
	DateCreated  time.Time `json:"date_created"`
	BucketPath   string    `json:"bucket_path"`
	CourseID     string    `json:"course_id"`
	InstructorID string    `json:"instructor_id"`
	SemesterTerm string    `json:"semester_term"`
	Section      string    `json:"section"`
}
