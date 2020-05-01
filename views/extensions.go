package views

import (
	"github.com/gorilla/mux"
	"html/template"
)

var router *mux.Router

func SetRouter(r *mux.Router) {
	router = r
}

func rawHTML(s string) template.HTML {
	return template.HTML(s)
}

func urlGenerator() func(string, ...string) string {
	return func(routeName string, pairs ...string) string {
		url, err := router.Get(routeName).URL(pairs...)
		if err != nil {
			panic(err)
		}

		return url.String()
	}
}
