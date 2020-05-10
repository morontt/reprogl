package handlers

import (
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"xelbot.com/reprogl/config"
	"xelbot.com/reprogl/models"
	"xelbot.com/reprogl/models/repositories"
	"xelbot.com/reprogl/views"
)

func PageAction(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		slug := vars["slug"]

		repo := repositories.ArticleRepository{DB: app.DB}
		article, err := repo.GetBySlug(slug)
		if err != nil {
			if errors.Is(err, models.RecordNotFound) {
				app.NotFound(w)
			} else {
				app.ServerError(w, err)
			}

			return
		}

		tagRepo := repositories.TagRepository{DB: app.DB}
		article.Tags, err = tagRepo.GetCollectionByArticle(article)
		if err != nil {
			app.ServerError(w, err)

			return
		}

		templateData := views.NewArticlePageData(article)
		templateData.AppendTitle(article.Title)

		err = views.RenderTemplate(w, "article.gohtml", templateData)
		if err != nil {
			app.ServerError(w, err)
		}
	}
}
