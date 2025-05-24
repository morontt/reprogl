package oauth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/oauth2"
	"xelbot.com/reprogl/container"
)

var vkEndpoint = oauth2.Endpoint{
	AuthURL:  "https://id.vk.com/authorize",
	TokenURL: "https://id.vk.com/oauth2/auth",
}

type vkResourceOwner struct {
	accessToken string
}

type vkUserInfoResponse struct {
	User vkUser `json:"user"`
}

type vkUser struct {
	ID        string `json:"user_id"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Gender    int    `json:"sex,omitempty"`
	Avatar    string `json:"avatar,omitempty"`
}

func (vkro *vkResourceOwner) GetUserData() (*UserData, error) {
	data := url.Values{}
	data.Set("access_token", vkro.accessToken)
	data.Set("client_id", container.GetConfig().OAuthVkID)

	request, err := http.NewRequest(
		http.MethodPost,
		"https://id.vk.com/oauth2/user_info",
		strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Accept", "application/json")

	buf, err := doRequest(request)
	if err != nil {
		return nil, err
	}

	if !json.Valid(buf) {
		return nil, errors.New("invalid JSON string")
	}

	result := vkUserInfoResponse{}
	err = json.Unmarshal(buf, &result)
	if err != nil {
		return nil, err
	}

	var vkGender Gender
	switch result.User.Gender {
	case 1:
		vkGender = Female
	case 2:
		vkGender = Male
	default:
		vkGender = Unknown
	}

	userData := &UserData{
		ID:          result.User.ID,
		DisplayName: strings.TrimSpace(result.User.FirstName + " " + result.User.LastName),
		FirstName:   result.User.FirstName,
		Gender:      vkGender,
		Avatar:      result.User.Avatar,
		RawData:     base64.URLEncoding.EncodeToString(buf),
		Provider:    vkProvider,
	}

	return userData, nil
}
