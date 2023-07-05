package models

type LoggedUser struct {
	ID           int
	Username     string
	PasswordHash string
	Salt         string
}
