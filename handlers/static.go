package handlers

import (
	"fmt"
	"net/http"
	"os"
	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/views"
)

func InfoAction(w http.ResponseWriter, _ *http.Request) {
	templateData := views.NewInfoPageData()
	templateData.AppendTitle("Информация")

	err := views.RenderTemplate(w, "info.gohtml", templateData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
	}
}

func RobotsTXTAction(w http.ResponseWriter, _ *http.Request) {
	var body string

	cfg := container.Get()
	body = fmt.Sprintf(
		"User-agent: *\n\nHost: %s\nSitemap: https://%s/sitemap.xml\n",
		cfg.Host,
		cfg.Host)

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(body))
}
