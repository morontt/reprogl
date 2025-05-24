package models

import (
	"database/sql"
	"strings"
	"time"

	"xelbot.com/reprogl/utils/hashid"
)

type LoggedUser struct {
	ID           int
	Username     string
	PasswordHash string
	Salt         string
	Role         string
}

type User struct {
	ID            int
	Username      string
	Email         string
	Role          string
	DisplayName   sql.NullString
	Gender        int
	AvatarVariant int
	CreatedAt     time.Time
}

func (u *User) Avatar(size int) string {
	options := hashid.User
	if u.Gender == 1 {
		options |= hashid.Male
	} else {
		options |= hashid.Female
	}

	if u.AvatarVariant > 0 {
		options += hashid.Option(u.AvatarVariant << 4)
	}

	return AvatarLink(u.ID, options, size)
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

func (u *User) DisplayNameValue() (name string) {
	if u.DisplayName.Valid && len(u.DisplayName.String) > 0 {
		name = u.DisplayName.String
	}

	return
}

func (u *User) IsMale() bool {
	return u.Gender == 1
}

func (u *User) NeedToCheckGravatar() bool {
	return u.HasEmail()
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) HasEmail() bool {
	return !strings.Contains(u.Email, "@xelbot.fake")
}
