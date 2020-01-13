package handlers

import (
	"fmt"
	"net/http"
	"os"
	"xelbot.com/reprogl/views"
)

func InfoAction(w http.ResponseWriter, _ *http.Request) {
	err := views.RenderTemplate(w, "static/info", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
	}
}

func RobotsTXTAction(w http.ResponseWriter, _ *http.Request) {
	var body string

	body = "User-agent: *\n\nHost: morontt.info\nSitemap: https://morontt.info/sitemap.xml\n"

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(body))
}
