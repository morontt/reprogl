package views

import (
	"crypto/sha256"
	"encoding/base64"
	"os"
	"xelbot.com/reprogl/container"
)

var fileHashes map[string]string

func init() {
	fileHashes = make(map[string]string)
}

func cdnBase() string {
	return cfg.CDNBaseURL
}

func subresourceIntegrity(file string) string {
	var n int

	h0, ok := fileHashes[file]
	if ok {
		return "sha256-" + h0
	}

	f, err := os.Open("./public/" + file)
	if err != nil {
		return ""
	}

	defer f.Close()

	hash := sha256.New()
	bytes := make([]byte, 1024)
	for {
		n, err = f.Read(bytes)
		if err != nil {
			break
		}

		hash.Write(bytes[:n])
	}

	h1 := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	if !container.IsDevMode() {
		fileHashes[file] = h1
	}

	return "sha256-" + h1
}
