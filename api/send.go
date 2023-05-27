package api

import (
	"net/http"
	"time"
)

func Send(req *http.Request) (*http.Response, error) {
	client := http.Client{
		Timeout: 15 * time.Second,
	}

	return client.Do(req)
}
