package backend

import (
	"bytes"
	"encoding/json"
	"net/http"

	"xelbot.com/reprogl/services/oauth"
)

type ExternalUserDTO struct {
	UserData  *oauth.UserData `json:"userData"`
	UserAgent string          `json:"userAgent"`
	IP        string          `json:"ipAddress"`
}

type CreatedUserDTO struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type CreateUserResponse struct {
	Violations []FormError     `json:"errors,omitempty"`
	User       *CreatedUserDTO `json:"comment,omitempty"`
}

func SendUserData(userData ExternalUserDTO) (*CreateUserResponse, error) {
	jsonBody, err := json.Marshal(userData)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, apiURL()+"/api/users/external", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")

	_, err = send(request)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
