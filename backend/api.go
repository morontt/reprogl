package backend

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/security"
)

type CommentatorDTO struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Website string `json:"website"`
}

type CommentDTO struct {
	Commentator CommentatorDTO `json:"commentator"`
	Text        string         `json:"text"`
	TopicID     int            `json:"topicId"`
	ParentID    int            `json:"parentId"`
	UserAgent   string         `json:"userAgent"`
	IP          string         `json:"ipAddress"`
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
