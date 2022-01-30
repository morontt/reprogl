package main

import (
	"github.com/gorilla/mux"
	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/handlers"
)

func getRoutes(app *container.Application) *mux.Router {
	siteMux := mux.NewRouter()
	siteMux.HandleFunc("/article/{slug}", handlers.PageAction(app)).Name("article")
	siteMux.HandleFunc("/{page:[0-9]*}", handlers.IndexAction(app)).Name("blog-page")
	siteMux.HandleFunc("/category/{slug}", handlers.CategoryAction(app)).Name("category-first")
	siteMux.HandleFunc("/category/{slug}/{page:[0-9]+}", handlers.CategoryAction(app)).Name("category")
	siteMux.HandleFunc("/tag/{slug}", handlers.TagAction(app)).Name("tag-first")
	siteMux.HandleFunc("/tag/{slug}/{page:[0-9]+}", handlers.TagAction(app)).Name("tag")
	siteMux.HandleFunc("/info", handlers.InfoAction).Name("info-page")
	siteMux.HandleFunc("/robots.txt", handlers.RobotsTXTAction)
	siteMux.HandleFunc("/favicon.ico", handlers.FavIconAction)

	fragmentsMux := siteMux.PathPrefix("/_fragment").Subrouter()
	fragmentsMux.HandleFunc("/categories", handlers.CategoriesFragment(app)).Name("fragment-categories")
	fragmentsMux.HandleFunc(
		"/comments/{article_id:[0-9]+}/{last_time:[0-9]+}",
		handlers.CommentsFragment(app),
	).Name("fragment-comments")

	app.Router = siteMux

	return siteMux
}
