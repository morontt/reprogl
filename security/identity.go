package security

import "xelbot.com/reprogl/models"

type Identity struct {
	ID       int
	Username string
}

func CreateIdentity(user *models.LoggedUser) Identity {
	return Identity{
		ID:       user.ID,
		Username: user.Username,
	}
}

func (i Identity) IsZero() bool {
	return i.ID == 0
}
