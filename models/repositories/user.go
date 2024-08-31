package repositories

import (
	"database/sql"
	"errors"
	"time"

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
			u.user_type,
			u.password,
			u.password_salt
		FROM users AS u
		WHERE (u.username = ?)`

	user := models.LoggedUser{}
	err := ur.DB.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Role,
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

func (ur *UserRepository) Find(id int) (*models.User, error) {
	query := `
		SELECT
			u.id,
			u.username,
			u.mail,
			u.user_type,
			u.display_name,
			u.avatar_variant,
			u.time_created,
			u.gender
		FROM users AS u
		WHERE (u.id = ?)`

	user := models.User{}
	err := ur.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Role,
		&user.DisplayName,
		&user.AvatarVariant,
		&user.CreatedAt,
		&user.Gender)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.RecordNotFound
		} else {
			return nil, err
		}
	}

	return &user, nil
}

func (ur *UserRepository) SaveLoginEvent(id int, ip string) error {
	query := `
		UPDATE
			users
		SET
			last_login = ?,
			login_count = login_count + 1,
			ip_last = ?
		WHERE
			id = ?`

	_, err := ur.DB.Exec(
		query,
		time.Now().Format("2006-01-02 15:04:05.000"),
		ip,
		id,
	)

	return err
}
