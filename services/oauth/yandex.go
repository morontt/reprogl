package oauth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
)

type yandexResourceOwner struct {
	accessToken string
}

func (yro *yandexResourceOwner) GetUserData() (*UserData, error) {
	request, err := http.NewRequest(http.MethodGet, "https://login.yandex.ru/info", http.NoBody)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "OAuth "+yro.accessToken)

	buf, err := doRequest(request)
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
