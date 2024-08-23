package backend

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"xelbot.com/reprogl/services/oauth"
)

type ExternalUserDTO struct {
	UserData  *oauth.UserData `json:"userData"`
	UserAgent string          `json:"userAgent"`
	IP        string          `json:"ipAddress"`
}

type CreatedUserDTO struct {
	ID          int    `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	DisplayName string `json:"displayName,omitempty"`
	ImageHash   string `json:"imageHash"`
}

type ProfileDTO struct {
	ID          int
	Username    string `json:"username"`
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
	Role        string `json:"role"`
	IsMale      bool   `json:"isMale"`
}

type CreateUserResponse struct {
	Violations []FormError     `json:"errors,omitempty"`
	User       *CreatedUserDTO `json:"user,omitempty"`
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

	response, err := send(request)
	if err != nil {
		return nil, err
	}

	if !(response.StatusCode == http.StatusOK ||
		response.StatusCode == http.StatusCreated ||
		response.StatusCode == http.StatusUnprocessableEntity) {
		return nil, errors.New("backend: unexpected HTTP status " + response.Status)
	}

	defer response.Body.Close()
	buf, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if !json.Valid(buf) {
		return nil, errors.New("invalid JSON string")
	}

	var result CreateUserResponse
	err = json.Unmarshal(buf, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func SendProfileData(profileData ProfileDTO) (*CreateUserResponse, error) {
	var requestData = struct {
		User ProfileDTO `json:"user"`
	}{
		User: profileData,
	}

	jsonBody, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPut, apiURL()+fmt.Sprintf("/api/users/%d", profileData.ID), bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := send(request)
	if err != nil {
		return nil, err
	}

	if !(response.StatusCode == http.StatusOK ||
		response.StatusCode == http.StatusUnprocessableEntity) {
		return nil, errors.New("backend: unexpected HTTP status " + response.Status)
	}

	defer response.Body.Close()
	buf, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if !json.Valid(buf) {
		return nil, errors.New("invalid JSON string")
	}

	var result CreateUserResponse
	err = json.Unmarshal(buf, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (u *CreatedUserDTO) Nickname() string {
	if len(u.DisplayName) > 0 {
		return u.DisplayName
	}

	return u.Username
}
