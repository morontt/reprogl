package backend

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/security"
)

type CommentDTO struct {
	Name     string
	Email    string
	Website  string
	Text     string
	TopicID  int
	ParentID int
}

var apiURL string

func init() {
	cnf := container.GetConfig()
	apiURL = cnf.BackendApiUrl
}

func SendComment(comment CommentDTO) error {
	jsonBody, err := json.Marshal(comment)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, apiURL+"/api/comments/external", bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}

	wsseHeader, wsseToken := security.GetWSSEHeader()

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set(wsseHeader, wsseToken)

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil
}
