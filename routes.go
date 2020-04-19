package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"xelbot.com/reprogl/handlers"
)

func getRoutes(app *handlers.Application) http.Handler {
	siteMux := mux.NewRouter()
	siteMux.HandleFunc("/article/{slug}", app.PageAction).Name("article")
	siteMux.HandleFunc("/{page:[0-9]*}", handlers.IndexAction).Name("blog-page")
	siteMux.HandleFunc("/category/{slug}/{page:[0-9]*}", handlers.CategoryAction).Name("category")
	siteMux.HandleFunc("/tag/{slug}/{page:[0-9]*}", handlers.TagAction).Name("tag")
	siteMux.HandleFunc("/info", handlers.InfoAction).Name("info_page")
	siteMux.HandleFunc("/robots.txt", handlers.RobotsTXTAction)

	return siteMux
}
