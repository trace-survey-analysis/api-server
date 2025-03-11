package database

import (
	"database/sql"
	"fmt"

	"api-server/internal/config"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB(cfg *config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func GetDB() *sql.DB {
	return db
}
