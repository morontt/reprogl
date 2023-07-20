package security

import (
	"errors"

	"xelbot.com/reprogl/models"
)

type Identity struct {
	ID       int    `json:"i,omitempty"`
	Username string `json:"u,omitempty"`
	Role     Role   `json:"r,omitempty"`
}

type IdentityAware interface {
	SetIdentity(Identity)
	HasIdentity() bool
	IsAdmin() bool
}

type Role uint8

const (
	Guest = iota
	User
	Admin
)

var (
	SerializeRoleError   = errors.New("security: serialize unknown role")
	DeserializeRoleError = errors.New("security: unserialize unknown role")
)

func CreateIdentity(user *models.LoggedUser) Identity {
	var role Role
	if user.Type == "admin" {
		role = Admin
	} else {
		role = User
	}

	return Identity{
		ID:       user.ID,
		Username: user.Username,
		Role:     role,
	}
}

func (i Identity) IsZero() bool {
	return i.ID == 0
}

func (i Identity) IsAdmin() bool {
	return i.Role == Admin
}

func (r Role) MarshalText() (text []byte, err error) {
	switch r {
	case Admin:
		text = []byte("admin")
	case User:
		text = []byte("user")
	case Guest:
		text = []byte{}
	default:
		err = SerializeRoleError
	}

	return
}

func (r *Role) UnmarshalText(text []byte) error {
	switch string(text) {
	case "admin":
		*r = Admin
	case "user":
		*r = User
	case "":
		*r = Guest
	default:
		return DeserializeRoleError
	}

	return nil
}
