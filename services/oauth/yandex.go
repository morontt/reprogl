package oauth

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"
)

type yandexResourceOwner struct {
	accessToken string
}

func (yro *yandexResourceOwner) GetUserData() (*UserData, error) {
	userData := &UserData{Provider: yandexProvider}

	request, err := http.NewRequest(http.MethodGet, "/api/comments/external", bytes.NewReader(nil))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Accept", "application/json")
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("oauth: response status code is " + strconv.Itoa(response.StatusCode))
	}

	defer response.Body.Close()
	buf, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if !json.Valid(buf) {
		return nil, errors.New("invalid JSON string")
	}

	result := struct {
		ID          string `json:"id"`
		Username    string `json:"login"`
		DisplayName string `json:"display_name"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Gender      string `json:"sex"` // check is NULL
		Email       string `json:"default_email"`
		Avatar      string `json:"default_avatar_id"`

		IsAvatarEmpty bool `json:"is_avatar_empty"`
	}{}

	err = json.Unmarshal(buf, &result)
	if err != nil {
		return nil, err
	}

	return userData, nil
}
