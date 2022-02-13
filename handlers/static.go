package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/views"
)

func InfoAction(w http.ResponseWriter, _ *http.Request) {
	templateData := views.NewInfoPageData()
	templateData.AppendTitle("Информация")

	err := views.WriteTemplate(w, "info.gohtml", templateData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
	}
}

func RobotsTXTAction(w http.ResponseWriter, _ *http.Request) {
	var body string

	cfg := container.GetConfig()
	body = fmt.Sprintf(
		"User-agent: *\n\nHost: %s\nSitemap: https://%s/sitemap.xml\n",
		cfg.Host,
		cfg.Host)

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(body))
}

func FavIconAction(w http.ResponseWriter, _ *http.Request) {
	var icon string
	icon = `
		AAABAAEAEBAQAAEABAAoAQAAFgAAACgAAAAQAAAAIAAAAAEABAAAAAAAAAAAAAAAAAAAAAAAAAAA
		AAAAAAAAAAAA////AACAAAAAgIAAgAAAAIAAgACAgAAAgICAAOPj4wAAAP8AAP8AAAD//wD/AAAA
		/wD/AP//tQD///8AAqKgAAAAoqIAKgiIiIgKKgAAiACIiICiAAiIiACIiAoAAAiIiIiIAgAe4IiI
		iIgAABHgiAAIiIAAERCAHuCIgAAACIAR4IiAAAiIgBEQiIAACAiIAAiIAAAAgIiIiIgAAAgIiIiI
		iAAAAICAgIiAAAAAAICAiAAAAAAAAAAAAACAAAAAwAAAAOAAAADAAAAAwAAAAIABAACAAAAAgAAA
		AMAAAADAAAAAwAEAAMABAADAAQAA4AMAAPAHAAD4DwAA`

	body, _ := base64.StdEncoding.DecodeString(icon)

	w.Header().Set("Content-Type", "image/x-icon")
	w.Header().Set("Cache-Control", "max-age=2592000")
	w.Write(body)
}
