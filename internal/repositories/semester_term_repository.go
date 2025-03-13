// get semester term by id
package repositories

import (
	"database/sql"

	"api-server/internal/models"
)

func GetSemesterTerm(db *sql.DB, semesterTerm string) (*models.SemesterTermModel, error) {
	semester := &models.SemesterTermModel{}
	err := db.QueryRow(
		"SELECT semester_term, name FROM webapp.semester_terms WHERE semester_term = $1",
		semesterTerm,
	).Scan(&semester.SemesterTerm, &semester.Name)
	return semester, err
}
