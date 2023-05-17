package backend

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
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

type ViolationPath string

type FormError struct {
	Message string        `json:"message"`
	Path    ViolationPath `json:"path"`
}

var apiURL string

func init() {
	cnf := container.GetConfig()
	apiURL = cnf.BackendApiUrl
}

func (vp *ViolationPath) UnmarshalText(jsonText []byte) error {
	text := string(jsonText)
	if strings.Index(text, "text") != -1 {
		*vp = "comment_text"
	} else if strings.Index(text, "email") != -1 {
		*vp = "mail"
	} else if strings.Index(text, "website") != -1 {
		*vp = "website"
	} else if strings.Index(text, "name") != -1 {
		*vp = "name"
	} else {
		*vp = "unknown"
	}

	return nil
}

func SendComment(comment CommentDTO) ([]FormError, error) {
	var violations []FormError
	jsonBody, err := json.Marshal(comment)
	if err != nil {
		return violations, err
	}

	request, err := http.NewRequest(http.MethodPost, apiURL+"/api/comments/external", bytes.NewReader(jsonBody))
	if err != nil {
		return violations, err
	}

	wsseHeader, wsseToken := security.GetWSSEHeader()

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set(wsseHeader, wsseToken)

	client := http.Client{
		Timeout: 15 * time.Second,
	}

	response, err := client.Do(request)
	if err != nil {
		return violations, err
	}

	if response.StatusCode == http.StatusUnprocessableEntity {
		buf, err := io.ReadAll(response.Body)
		if err != nil {
			return violations, err
		}

		if !json.Valid(buf) {
			return violations, errors.New("invalid JSON string")
		}

		err = json.Unmarshal(buf, &violations)
		if err != nil {
			return violations, err
		}
	}

	return violations, nil
}
