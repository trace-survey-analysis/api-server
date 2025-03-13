package models

type Department struct {
	DepartmentID int    `json:"department_id"`
	Name         string `json:"name"`
	SchoolID     int    `json:"school_id"`
}
