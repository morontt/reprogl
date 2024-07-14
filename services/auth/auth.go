package auth

import (
	"crypto/subtle"
	"errors"
	"fmt"

	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/models"
	"xelbot.com/reprogl/models/repositories"
	"xelbot.com/reprogl/security"
)

type Error interface {
	InfoLogMessage() string
}

type wrongCredentialsError struct {
	info string
}

func HandleLoginPassword(app *container.Application, username, password string) (*models.LoggedUser, error) {
	if len(username) == 0 || len(password) == 0 {
		return nil, wrongCredentialsError{info: "username or password is empty"}
	}

	if len(password) > 4096 {
		return nil, wrongCredentialsError{info: "password too long"}
	}

	repo := repositories.UserRepository{DB: app.DB}
	user, err := repo.GetLoggedUserByUsername(username)
	if err != nil {
		if errors.Is(err, models.RecordNotFound) {
			return nil, wrongCredentialsError{info: fmt.Sprintf("user \"%s\" not found", username)}
		}

		return nil, err
	}

	passwordHash := security.EncodePassword(password, user.Salt)
	if subtle.ConstantTimeCompare([]byte(passwordHash), []byte(user.PasswordHash)) == 0 {
		return nil, wrongCredentialsError{info: fmt.Sprintf("invalid password for \"%s\"", username)}
	}

	return user, nil
}

func (_ wrongCredentialsError) Error() string {
	return "Недействительные логин/пароль"
}

func (w wrongCredentialsError) InfoLogMessage() string {
	return w.info
}
