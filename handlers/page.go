package handlers

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"xelbot.com/reprogl/models"
	"xelbot.com/reprogl/models/repositories"
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
			app.serverError(w, err)
		}

		return
	}

	fmt.Fprintf(w, "Article %s", article.Text)
}
