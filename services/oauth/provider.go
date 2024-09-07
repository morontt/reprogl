package oauth

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/yandex"
	"xelbot.com/reprogl/container"
)

var (
	ProviderNotFound = errors.New("oauth: no matching provider found")
	logger           *log.Logger
)

const (
	yandexProvider = "yandex"
	vkProvider     = "vkontakte"
)

func SetLogger(l *log.Logger) {
	logger = l
}

func SupportedProvider(name string) (found bool) {
	switch name {
	case yandexProvider:
		found = true
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
	case vkProvider:
		return &oauth2.Config{
			ClientID:     cnf.OAuthVkID,
			ClientSecret: cnf.OAuthVkSecret,
			Endpoint:     vkEndpoint,
			RedirectURL:  url,
			Scopes:       []string{"vkid.personal_info", "email"},
		}, nil
	}

	return nil, ProviderNotFound
}

func AdditionalParams(name string) []string {
	params := make([]string, 0)
	switch name {
	case vkProvider:
		return []string{"device_id"}
	}

	return params
}

func UserDataByCode(providerName, code, verifier string, additional map[string]string) (*UserData, error) {
	oauthConfig, err := ConfigByProvider(providerName)
	if err != nil {
		return nil, err
	}

	options := make([]oauth2.AuthCodeOption, 1)
	options[0] = oauth2.VerifierOption(verifier)

	for key, value := range additional {
		options = append(options, oauth2.SetAuthURLParam(key, value))
	}

	token, err := oauthConfig.Exchange(context.Background(), code, options...)
	if err != nil {
		return nil, err
	}

	if tokenScope := token.Extra("scope"); tokenScope != nil {
		logger.Printf("[OAUTH] token scope: %s", tokenScope)
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
	case vkProvider:
		return &vkResourceOwner{accessToken: token.AccessToken}, nil
	}

	return nil, ProviderNotFound
}

func doRequest(request *http.Request) ([]byte, error) {
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

	logger.Printf("[OAUTH] user data: %s", buf)

	return buf, nil
}
