package oauth

import (
	"errors"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/yandex"
	"xelbot.com/reprogl/container"
)

func ConfigByProvider(name string) (*oauth2.Config, error) {
	cnf := container.GetConfig()

	switch name {
	case "yandex":
		return &oauth2.Config{
			ClientID:     cnf.OAuthYandexID,
			ClientSecret: cnf.OAuthYandexSecret,
			Endpoint:     yandex.Endpoint,
		}, nil
	}

	return nil, errors.New("oauth: no matching provider found")
}
