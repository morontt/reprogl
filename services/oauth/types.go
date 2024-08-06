package oauth

import (
	"encoding/json"
	"strings"
)

type Gender string

type ResourceOwnerInterface interface {
	GetUserData() (*UserData, error)
}

const (
	Male    Gender = "male"
	Female  Gender = "female"
	Unknown Gender = "n/a"
)

type UserData struct {
	ID          string
	Username    string
	DisplayName string
	FirstName   string
	LastName    string
	Gender      Gender
	Email       string
	Avatar      string
	Provider    string

	RawData string
}

func (a *Gender) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	switch strings.ToLower(s) {
	case "male":
		*a = Male
	case "female":
		*a = Female
	default:
		*a = Unknown
	}

	return nil
}
