package oauth

import "golang.org/x/oauth2"

var vkEndpoint = oauth2.Endpoint{
	AuthURL:  "https://id.vk.com/authorize",
	TokenURL: "https://id.vk.com/oauth2/auth",
}
