// get semester term by id
package repositories

import (
	"database/sql"

	"api-server/internal/models"
)

func GetSemesterTerm(db *sql.DB, semesterTerm string) (*models.SemesterTermModel, error) {
	semester := &models.SemesterTermModel{}
	err := db.QueryRow(
		"SELECT semester_term, name FROM api.semester_terms WHERE semester_term = $1",
		semesterTerm,
	).Scan(&semester.SemesterTerm, &semester.Name)
	return semester, err
}

// GetAllSemesterTerms retrieves all semester terms
func GetAllSemesterTerms(db *sql.DB) ([]models.SemesterTermModel, error) {
	rows, err := db.Query(
		"SELECT semester_term, name FROM api.semester_terms",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	semesterTerms := []models.SemesterTermModel{}
	for rows.Next() {
		var semesterTerm models.SemesterTermModel
		if err := rows.Scan(&semesterTerm.SemesterTerm, &semesterTerm.Name); err != nil {
			return nil, err
		}
		semesterTerms = append(semesterTerms, semesterTerm)
	}
	return semesterTerms, nil
}
