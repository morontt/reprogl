package backend

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"xelbot.com/reprogl/api"
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

type CreatedCommentDTO struct {
	ID      int    `json:"id"`
	Text    string `json:"text"`
	Name    string `json:"username"`
	Email   string `json:"email"`
	Website string `json:"website"`
	Country string `json:"countryCode"`
}

type ViolationPath string

type FormError struct {
	Message string        `json:"message"`
	Path    ViolationPath `json:"path"`
}

type CreateCommentResponse struct {
	Violations []FormError        `json:"errors,omitempty"`
	Comment    *CreatedCommentDTO `json:"comment,omitempty"`
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

func SendComment(comment CommentDTO) (*CreateCommentResponse, error) {
	var result CreateCommentResponse
	jsonBody, err := json.Marshal(comment)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, apiURL+"/api/comments/external", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	wsseHeader, wsseToken := security.GetWSSEHeader()

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set(wsseHeader, wsseToken)

	response, err := api.Send(request)
	if err != nil {
		return nil, err
	}

	buf, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if !json.Valid(buf) {
		return nil, errors.New("invalid JSON string")
	}

	err = json.Unmarshal(buf, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
