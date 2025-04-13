package repositories

import (
	"api-server/internal/models"
	"database/sql"
	"time"
)

// CreateCourse creates a new course in the database

func CreateCourse(db *sql.DB, course models.Course) (models.Course, error) {

	_, err := db.Exec(
		"INSERT INTO api.courses (course_id, date_added, date_last_updated, user_id, code, name, description, instructor_id, department_id, credit_hours) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
		course.CourseID, course.DateAdded, course.DateLastUpdated, course.UserID, course.Code, course.Name, course.Description, course.InstructorID, course.DepartmentID, course.CreditHours,
	)
	if err != nil {
		return models.Course{}, err
	}
	return course, err
}

// GetCourseByID retrieves a course by its ID
func GetCourseByID(db *sql.DB, courseID string) (*models.Course, error) {
	course := &models.Course{}
	err := db.QueryRow(
		"SELECT course_id, date_added, date_last_updated, user_id, code, name, description, instructor_id, department_id, credit_hours FROM api.courses WHERE course_id = $1",
		courseID,
	).Scan(&course.CourseID, &course.DateAdded, &course.DateLastUpdated, &course.UserID, &course.Code, &course.Name, &course.Description, &course.InstructorID, &course.DepartmentID, &course.CreditHours)
	return course, err
}

// GetAllCourses retrieves all courses
func GetAllCourses(db *sql.DB) ([]models.Course, error) {
	rows, err := db.Query(
		"SELECT course_id, date_added, date_last_updated, user_id, code, name, description, instructor_id, department_id, credit_hours FROM api.courses",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	courses := []models.Course{}
	for rows.Next() {
		var course models.Course
		if err := rows.Scan(&course.CourseID, &course.DateAdded, &course.DateLastUpdated, &course.UserID, &course.Code, &course.Name, &course.Description, &course.InstructorID, &course.DepartmentID, &course.CreditHours); err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}
	return courses, nil
}

// UpdateCourse updates a course in the database
func UpdateCourse(db *sql.DB, course *models.Course) error {
	course.DateLastUpdated = time.Now().UTC()
	_, err := db.Exec(
		"UPDATE api.courses SET date_last_updated=$1, code=$2, name=$3, description=$4, instructor_id=$5, department_id=$6, credit_hours=$7 WHERE course_id=$8",
		course.DateLastUpdated, course.Code, course.Name, course.Description, course.InstructorID, course.DepartmentID, course.CreditHours, course.CourseID,
	)
	return err
}

// delete course by ID
func DeleteCourse(db *sql.DB, courseID string) error {
	result, err := db.Exec(
		"DELETE FROM api.courses WHERE course_id = $1",
		courseID,
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
