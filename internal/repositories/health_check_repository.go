package repositories

import "database/sql"

// InsertHealthCheck inserts a new health check record into the database.
func InsertHealthCheck(db *sql.DB) error {
	_, err := db.Exec("INSERT INTO webapp.health_checks (checked_at) VALUES (CURRENT_TIMESTAMP)")
	return err
}
