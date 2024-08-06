package avatar

import (
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"runtime"
	"strings"
	"time"

	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/models"
)

var (
	GravatarNotFound = errors.New("gravatar: not found")
	chanBuf          chan struct{}
)

func init() {
	chanBuf = make(chan struct{}, (runtime.NumCPU()/2)+1)
}

func tryGravatar(commentator *models.CommentatorForGravatar) (image.Image, error) {
	if !commentator.Email.Valid || !commentator.FakeEmail.Valid || commentator.FakeEmail.Bool {
		return nil, GravatarNotFound
	}

	request, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(
			"https://www.gravatar.com/avatar/%s?s=80&d=404",
			gravatarHash(commentator.Email.String),
		),
		nil)

	if err != nil {
		return nil, err
	}

	resp, err := sendRequest(request)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, GravatarNotFound
	}

	contentType, ok := resp.Header["Content-Type"]
	if !ok || !strings.HasPrefix(contentType[0], "image") {
		return nil, GravatarNotFound
	}

	defer resp.Body.Close()
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func gravatarHash(s string) string {
	return container.MD5(strings.ToLower(strings.TrimSpace(s)))
}

func sendRequest(req *http.Request) (*http.Response, error) {
	chanBuf <- struct{}{}
	defer func() { <-chanBuf }()

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11")
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	return client.Do(req)
}
