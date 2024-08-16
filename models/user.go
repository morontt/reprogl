package models

import (
	"database/sql"
	"time"

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
	CreatedAt   time.Time
}

func (u *User) Avatar(size int) string {
	options := hashid.User
	if u.Gender == 1 {
		options |= hashid.Male
	} else {
		options |= hashid.Female
	}

	return avatarLink(u.ID, options, size)
}

func (u *User) Nickname() string {
	if u.DisplayName.Valid && len(u.DisplayName.String) > 0 {
		return u.DisplayName.String
	}

	return u.Username
}

func (u *User) HasDisplayName() bool {
	return u.DisplayName.Valid && len(u.DisplayName.String) > 0
}
