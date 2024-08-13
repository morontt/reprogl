package models

import (
	"database/sql"

	"xelbot.com/reprogl/utils/hashid"
)

type LoggedUser struct {
	ID           int
	Username     string
	PasswordHash string
	Salt         string
	Type         string
}

type User struct {
	ID          int
	Username    string
	Email       string
	Type        string
	DisplayName sql.NullString
	Gender      int
}

func (u *User) Avatar() string {
	options := hashid.User
	if u.Gender == 1 {
		options |= hashid.Male
	} else {
		options |= hashid.Female
	}

	return avatarLink(u.ID, options)
}
