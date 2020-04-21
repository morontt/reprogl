package handlers

import (
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"xelbot.com/reprogl/models"
	"xelbot.com/reprogl/models/repositories"
	"xelbot.com/reprogl/views"
)

func (app *Application) PageAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	repo := repositories.ArticleRepository{DB: app.DB}
	article, err := repo.GetBySlug(slug)
	if err != nil {
		if errors.Is(err, models.RecordNotFound) {
			app.notFound(w)
		} else {
			app.ServerError(w, err)
		}

		return
	}

	templateData := views.ArticlePageData{Article: article}

	err = views.RenderTemplate(w, "article.gohtml", templateData)
	if err != nil {
		app.ServerError(w, err)
	}
}
