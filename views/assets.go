package views

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"html/template"
	"os"
	"strconv"
	"strings"

	"xelbot.com/reprogl/container"
)

var fileHashes map[string]string

func init() {
	fileHashes = make(map[string]string)
}

func cdnBase() string {
	return cfg.CDNBaseURL
}

func assetTag(file string) template.HTML {
	var (
		tmpl    string
		result  string
		subPath string
		ext     string
	)

	result, ok := fileHashes[file]
	if ok {
		return template.HTML(result)
	}

	if strings.HasSuffix(file, ".css") {
		tmpl = `<link rel="stylesheet" href="%s" integrity="%s" crossorigin="anonymous">`
		subPath = file[:len(file)-3]
		ext = ".css"
	} else if strings.HasSuffix(file, ".js") {
		tmpl = `<script src="%s" integrity="%s" crossorigin="anonymous"></script>`
		subPath = file[:len(file)-2]
		ext = ".js"
	} else {
		tmpl = `<unknown src="%s" integrity="%s"/>`
		subPath = file
		ext = ".?"
	}

	hash := subresourceIntegrity(file)
	timeBasedPart := timeBased(file)

	result = fmt.Sprintf(tmpl, cfg.CDNBaseURL+"/"+subPath+timeBasedPart+ext, hash)
	if !container.IsDevMode() {
		fileHashes[file] = result
	}

	return template.HTML(result)
}

func subresourceIntegrity(file string) string {
	var n int

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

	return "sha256-" + base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func timeBased(file string) string {
	fileInfo, err := os.Stat("./public/" + file)
	if err != nil {
		return "v999"
	}

	modificatedAt := fileInfo.ModTime()

	return fmt.Sprintf(
		"v%d%s%s%s%s",
		modificatedAt.Year()%100,
		zeroPad(modificatedAt.YearDay(), 3),
		zeroPad(modificatedAt.Hour(), 2),
		zeroPad(modificatedAt.Minute(), 2),
		zeroPad(modificatedAt.Second(), 2),
	)
}

func zeroPad(i, width int) string {
	istring := strconv.Itoa(i)

	return strings.Repeat("0", width-len(istring)) + istring
}
