package controllers

import (
	"github.com/CloudyKit/jet"
	"net/http"
	"reprogl/views"
)

func InfoAction(w http.ResponseWriter, _ *http.Request) {
	tmpl, err := views.ViewSet.GetTemplate("static/info.jet")
	if err != nil {
		panic(err)
	}

	vars := make(jet.VarMap)
	tmpl.Execute(w, vars, nil)
}

func RobotsTXTAction(w http.ResponseWriter, _ *http.Request) {
	var body string

	body = "User-agent: *\n\nHost: morontt.info\nSitemap: https://morontt.info/sitemap.xml\n"

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(body))
}
