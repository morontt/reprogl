package backend

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"xelbot.com/reprogl/api"
	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/security"
)

type CommentatorDTO struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Website string `json:"website"`
}

type UserDTO struct {
	ID int `json:"id"`
}

type CommentDTO struct {
	Commentator *CommentatorDTO `json:"commentator,omitempty"`
	User        *UserDTO        `json:"user,omitempty"`
	Text        string          `json:"text"`
	TopicID     int             `json:"topicId"`
	ParentID    int             `json:"parentId"`
	UserAgent   string          `json:"userAgent"`
	IP          string          `json:"ipAddress"`
}

type CreatedCommentDTO struct {
	ID      int    `json:"id"`
	Text    string `json:"text"`
	Name    string `json:"username"`
	Email   string `json:"email"`
	Website string `json:"website"`
	Country string `json:"countryFlag"`
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

var backendLocker sync.Mutex

var NotAllowedComment = errors.New("backend: not allowed comment")

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

	request, err := http.NewRequest(http.MethodPost, apiURL()+"/api/comments/external", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := send(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode == http.StatusForbidden {
		return nil, NotAllowedComment
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

func PingGeolocation() {
	request, err := http.NewRequest(http.MethodPost, apiURL()+"/api/comments/geo-location", bytes.NewReader(nil))
	if err != nil {
		return
	}

	_, _ = send(request)
}

func RefreshComment(id int) (*CreatedCommentDTO, error) {
	request, err := http.NewRequest(http.MethodGet, apiURL()+"/api/comments/"+strconv.Itoa(id), bytes.NewReader(nil))
	if err != nil {
		return nil, err
	}

	response, err := send(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("status not OK")
	}

	buf, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if !json.Valid(buf) {
		return nil, errors.New("invalid JSON string")
	}

	var result = struct {
		Comment *CreatedCommentDTO `json:"comment,omitempty"`
	}{}
	err = json.Unmarshal(buf, &result)
	if err != nil {
		return nil, err
	}

	return result.Comment, nil
}

func send(req *http.Request) (*http.Response, error) {
	backendLocker.Lock()
	defer backendLocker.Unlock()

	wsseHeader, wsseToken := security.GetWSSEHeader()
	req.Header.Set(wsseHeader, wsseToken)

	return api.Send(req)
}

func apiURL() string {
	return container.GetConfig().BackendApiUrl
}
