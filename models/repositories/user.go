package repositories

import (
	"database/sql"
	"errors"

	"xelbot.com/reprogl/models"
)

type UserRepository struct {
	DB *sql.DB
}

func (ur *UserRepository) GetLoggedUserByUsername(username string) (*models.LoggedUser, error) {
	query := `
		SELECT
			u.id,
			u.username,
			u.password,
			u.password_salt
		FROM users AS u
		WHERE (u.username = ?)`

	user := models.LoggedUser{}

	err := ur.DB.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Salt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.RecordNotFound
		} else {
			return nil, err
		}
	}

	return &user, nil
}
