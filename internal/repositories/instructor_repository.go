package repositories

import (
	"api-server/internal/models"
	"database/sql"
)

// CreateInstructor inserts a new instructor record into the database.
func CreateInstructor(db *sql.DB, instructor models.Instructor) (models.Instructor, error) {
	query := `
        INSERT INTO webapp.instructors (instructor_id, user_id, name, date_created)
        VALUES ($1, $2, $3, $4)
    `
	_, err := db.Exec(query, instructor.InstructorID, instructor.UserID, instructor.Name, instructor.DateCreated)
	if err != nil {
		return models.Instructor{}, err
	}
	return instructor, nil
}

// GetInstructorByID retrieves an instructor by instructor_id.
func GetInstructorByID(db *sql.DB, instructorID string) (models.Instructor, error) {
	query := `
        SELECT instructor_id, user_id, name, date_created
        FROM webapp.instructors
        WHERE instructor_id = $1
    `
	row := db.QueryRow(query, instructorID)
	var instructor models.Instructor
	err := row.Scan(&instructor.InstructorID, &instructor.UserID, &instructor.Name, &instructor.DateCreated)
	if err != nil {
		return models.Instructor{}, err
	}
	return instructor, nil
}

// UpdateInstructor updates the instructor's name.
func UpdateInstructor(db *sql.DB, instructor models.Instructor) error {
	query := `
        UPDATE webapp.instructors
        SET name = $1
        WHERE instructor_id = $2
    `
	_, err := db.Exec(query, instructor.Name, instructor.InstructorID)
	return err
}

// DeleteInstructor deletes an instructor by instructor_id.
func DeleteInstructor(db *sql.DB, instructorID string) error {
	query := `DELETE FROM webapp.instructors WHERE instructor_id = $1`
	result, err := db.Exec(query, instructorID)
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
