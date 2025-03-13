package repositories

import (
	"database/sql"

	"api-server/internal/models"
)

// CreateTrace creates a new trace in the database
func CreateTrace(db *sql.DB, trace models.Trace) (models.Trace, error) {
	_, err := db.Exec(
		"INSERT INTO webapp.traces (trace_id, user_id, file_name, date_created, bucket_path, course_id, instructor_id, semester_term, section) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
		trace.TraceID, trace.UserID, trace.FileName, trace.DateCreated, trace.BucketPath, trace.CourseID, trace.InstructorID, trace.SemesterTerm, trace.Section,
	)
	if err != nil {
		return models.Trace{}, err
	}
	return trace, err
}

// GetTraceByID retrieves a trace by its ID
func GetTraceByID(db *sql.DB, traceID string) (*models.Trace, error) {
	trace := &models.Trace{}
	err := db.QueryRow(
		"SELECT trace_id, user_id, file_name, date_created, bucket_path, course_id, instructor_id, semester_term, section FROM webapp.traces WHERE trace_id = $1",
		traceID,
	).Scan(&trace.TraceID, &trace.UserID, &trace.FileName, &trace.DateCreated, &trace.BucketPath, &trace.CourseID, &trace.InstructorID, &trace.SemesterTerm, &trace.Section)
	return trace, err
}

// get all trace by courseID
func GetTraceByCourseID(db *sql.DB, courseID string) ([]models.Trace, error) {
	rows, err := db.Query(
		"SELECT trace_id, user_id, file_name, date_created, bucket_path, course_id, instructor_id, semester_term, section FROM webapp.traces WHERE course_id = $1",
		courseID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	traces := []models.Trace{}
	for rows.Next() {
		var trace models.Trace
		if err := rows.Scan(&trace.TraceID, &trace.UserID, &trace.FileName, &trace.DateCreated, &trace.BucketPath, &trace.CourseID, &trace.InstructorID, &trace.SemesterTerm, &trace.Section); err != nil {
			return nil, err
		}
		traces = append(traces, trace)
	}
	return traces, nil
}

// delete trace by ID
func DeleteTrace(db *sql.DB, traceID string) error {
	result, err := db.Exec(
		"DELETE FROM webapp.traces WHERE trace_id = $1",
		traceID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// get filepath from trace id
func GetFilePath(db *sql.DB, traceID string) (string, error) {
	var filePath string
	err := db.QueryRow(
		"SELECT bucket_path FROM webapp.traces WHERE trace_id = $1",
		traceID,
	).Scan(&filePath)
	return filePath, err
}
