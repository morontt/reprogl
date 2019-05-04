package controllers

import (
	"fmt"
	"net/http"
)

func InfoAction(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "Info page")
}

func RobotsTXTAction(w http.ResponseWriter, _ *http.Request) {
	var body string

	body = "User-agent: *\n\nHost: morontt.info\nSitemap: https://morontt.info/sitemap.xml\n"

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(body))
}
