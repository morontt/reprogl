package oauth

import (
	"context"
	"errors"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/yandex"
	"xelbot.com/reprogl/container"
)

var ProviderNotFound = errors.New("oauth: no matching provider found")

const (
	yandexProvider = "yandex"
	vkProvider     = "vkontakte"
)

func SupportedProvider(name string) (found bool) {
	switch name {
	case yandexProvider:
	case vkProvider:
		found = true
	}

	return
}

func ConfigByProvider(name string) (*oauth2.Config, error) {
	cnf := container.GetConfig()
	url := container.GenerateAbsoluteURL("oauth-verification", "provider", name)

	switch name {
	case yandexProvider:
		return &oauth2.Config{
			ClientID:     cnf.OAuthYandexID,
			ClientSecret: cnf.OAuthYandexSecret,
			Endpoint:     yandex.Endpoint,
			RedirectURL:  url,
		}, nil
	case vkProvider: // invalid code_challenge
		return &oauth2.Config{
			ClientID:     cnf.OAuthVkID,
			ClientSecret: cnf.OAuthVkSecret,
			Endpoint:     vkEndpoint,
			RedirectURL:  url,
			Scopes:       []string{"email"},
		}, nil
	}

	return nil, ProviderNotFound
}

func UserDataByCode(providerName, code string) (*UserData, error) {
	oauthConfig, err := ConfigByProvider(providerName)
	if err != nil {
		return nil, err
	}

	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}

	resourceOwner, err := resourceOwnerByProvider(providerName, token)
	if err != nil {
		return nil, err
	}

	return resourceOwner.GetUserData()
}

func resourceOwnerByProvider(name string, token *oauth2.Token) (ResourceOwnerInterface, error) {
	switch name {
	case yandexProvider:
		return &yandexResourceOwner{accessToken: token.AccessToken}, nil
	}

	return nil, ProviderNotFound
}
