package backend

import (
	"net/http"
	"sync"

	"xelbot.com/reprogl/api"
	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/security"
)

type FormError struct {
	Message string `json:"message"`
	Path    string `json:"path"`
}

var backendLocker sync.Mutex

func send(req *http.Request) (*http.Response, error) {
	backendLocker.Lock()
	defer backendLocker.Unlock()

	req.Header.Set("Authorization", "WSSE profile=\"UsernameToken\"")

	wsseHeader, wsseToken := security.GetWSSEHeader()
	req.Header.Set(wsseHeader, wsseToken)

	return api.Send(req)
}

func apiURL() string {
	return container.GetConfig().BackendApiUrl
}
