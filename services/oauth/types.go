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
	ID          string `json:"id"`
	Username    string `json:"username,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	FirstName   string `json:"firstName,omitempty"`
	LastName    string `json:"lastName,omitempty"`
	Gender      Gender `json:"gender,omitempty"`
	Email       string `json:"email,omitempty"`
	Avatar      string `json:"avatar,omitempty"`
	Provider    string `json:"dataProvider"`

	RawData string `json:"rawData"`
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
