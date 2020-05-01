package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"xelbot.com/reprogl/models/repositories"
	"xelbot.com/reprogl/views"
)

func (app *Application) IndexAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	page, needsRedirect := pageOrRedirect(vars)
	if needsRedirect {
		http.Redirect(w, r, "/", 301)

		return
	}

	repo := repositories.ArticleRepository{DB: app.DB}
	articles, err := repo.GetCollection(page)
	if err != nil {
		app.ServerError(w, err)

		return
	}

	templateData := views.IndexPageData{PageNumber: page, Articles: articles}

	err = views.RenderTemplate(w, "index.gohtml", templateData)
	if err != nil {
		app.ServerError(w, err)
	}
}

func CategoryAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	page, needsRedirect := pageOrRedirect(vars)
	if needsRedirect {
		http.Redirect(w, r, "/", 301)

		return
	}

	fmt.Fprintf(w, "Articles by category, page %d", page)
}

func TagAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	page, needsRedirect := pageOrRedirect(vars)
	if needsRedirect {
		http.Redirect(w, r, "/", 301)

		return
	}

	fmt.Fprintf(w, "Articles by tag, page %d", page)
}
