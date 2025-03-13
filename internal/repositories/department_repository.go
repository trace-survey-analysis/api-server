// get department by id
package repositories

import (
	"database/sql"

	"api-server/internal/models"
)

func GetDepartmentByID(db *sql.DB, departmentID int) (*models.Department, error) {
	department := &models.Department{}
	err := db.QueryRow(
		"SELECT department_id, name FROM api.departments WHERE department_id = $1",
		departmentID,
	).Scan(&department.DepartmentID, &department.Name)
	return department, err
}
