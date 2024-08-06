package oauth

import (
	"bytes"
	"encoding/base64"
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
	request, err := http.NewRequest(http.MethodGet, "https://login.yandex.ru/info", bytes.NewReader(nil))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "OAuth "+yro.accessToken)
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
		ID            string `json:"id"`
		Username      string `json:"login"`
		DisplayName   string `json:"display_name"`
		FirstName     string `json:"first_name"`
		LastName      string `json:"last_name"`
		Gender        Gender `json:"sex"`
		Email         string `json:"default_email"`
		Avatar        string `json:"default_avatar_id"`
		IsAvatarEmpty bool   `json:"is_avatar_empty"`
	}{
		Gender: Unknown,
	}

	err = json.Unmarshal(buf, &result)
	if err != nil {
		return nil, err
	}

	userData := &UserData{
		ID:          result.ID,
		Username:    result.Username,
		DisplayName: result.DisplayName,
		FirstName:   result.FirstName,
		LastName:    result.LastName,
		Gender:      result.Gender,
		Email:       result.Email,
		RawData:     base64.URLEncoding.EncodeToString(buf),
		Provider:    yandexProvider,
	}

	if !result.IsAvatarEmpty {
		userData.Avatar = result.Avatar
	}

	return userData, nil
}
