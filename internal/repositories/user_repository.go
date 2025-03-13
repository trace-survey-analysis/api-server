package repositories

import (
	"api-server/internal/models"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

func CreateUser(db *sql.DB, userReq models.UserRequest) (*models.User, error) {
	userID := uuid.New().String()
	now := time.Now().UTC()

	user := models.User{
		UserID:         userID,
		FirstName:      userReq.FirstName,
		LastName:       userReq.LastName,
		Password:       userReq.Password,
		Username:       userReq.Username,
		AccountCreated: now,
		AccountUpdated: now,
	}

	_, err := db.Exec(
		"INSERT INTO webapp.users (user_id, first_name, last_name, username, password, account_created, account_updated) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		user.UserID, user.FirstName, user.LastName, user.Username, user.Password, user.AccountCreated, user.AccountUpdated,
	)
	return &user, err
}

func GetUserByID(db *sql.DB, userID string) (*models.User, error) {
	user := &models.User{}
	err := db.QueryRow(
		"SELECT user_id, first_name, last_name, username, password, account_created, account_updated FROM webapp.users WHERE user_id = $1",
		userID,
	).Scan(&user.UserID, &user.FirstName, &user.LastName, &user.Username, &user.Password, &user.AccountCreated, &user.AccountUpdated)
	return user, err
}

func UpdateUser(db *sql.DB, user *models.User) error {
	user.AccountUpdated = time.Now().UTC()
	_, err := db.Exec(
		"UPDATE webapp.users SET first_name=$1, last_name=$2, username=$3, password=$4, account_updated=$5 WHERE user_id=$6",
		user.FirstName, user.LastName, user.Username, user.Password, user.AccountUpdated, user.UserID,
	)
	return err
}

func GetUserByUsername(db *sql.DB, username string) (*models.User, error) {
	user := &models.User{}
	err := db.QueryRow(
		"SELECT user_id, first_name, last_name, username FROM webapp.users WHERE username = $1",
		username,
	).Scan(&user.UserID, &user.FirstName, &user.LastName, &user.Username)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// UserWithPassword includes password for auth purposes
type UserWithPassword struct {
	models.User
	Password string
}

// GetUserWithPasswordByUsername retrieves a user with password by username
func GetUserWithPasswordByUsername(db *sql.DB, username string) (*UserWithPassword, error) {
	user := &UserWithPassword{}
	err := db.QueryRow(
		"SELECT user_id, first_name, last_name, username, password, account_created, account_updated FROM webapp.users WHERE username = $1",
		username,
	).Scan(&user.UserID, &user.FirstName, &user.LastName, &user.Username, &user.Password, &user.AccountCreated, &user.AccountUpdated)

	if err != nil {
		return nil, err
	}
	return user, nil
}
