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

// GetAllDepartments retrieves all departments
func GetAllDepartments(db *sql.DB) ([]models.Department, error) {
	rows, err := db.Query(
		"SELECT department_id, name FROM api.departments",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	departments := []models.Department{}
	for rows.Next() {
		var department models.Department
		if err := rows.Scan(&department.DepartmentID, &department.Name); err != nil {
			return nil, err
		}
		departments = append(departments, department)
	}
	return departments, nil
}
