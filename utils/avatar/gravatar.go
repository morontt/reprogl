package avatar

import (
	"errors"
	"fmt"
	"image"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"

	"xelbot.com/reprogl/container"
)

var (
	GravatarNotFound = errors.New("gravatar: not found")
	chanBuf          chan struct{}
)

func init() {
	chanBuf = make(chan struct{}, (runtime.NumCPU()/2)+1)
}

type MaybeGravatar interface {
	NeedToCheckGravatar() bool
	GetEmail() string
}

func tryGravatar(object MaybeGravatar, size int, logger *log.Logger) (image.Image, error) {
	if !object.NeedToCheckGravatar() {
		return nil, GravatarNotFound
	}

	hash := gravatarHash(object.GetEmail())
	request, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(
			"https://www.gravatar.com/avatar/%s?s=%d&d=404",
			hash,
			size,
		),
		nil)

	if err != nil {
		return nil, err
	}

	logger.Printf("[IMG] check gravatar hash %s on gravatar.com\n", hash)
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
