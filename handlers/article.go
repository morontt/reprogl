package handlers

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/models"
	"xelbot.com/reprogl/models/repositories"
	"xelbot.com/reprogl/views"
)

func PageAction(app *container.Application) http.HandlerFunc {
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

		commentRepo := repositories.CommentRepository{DB: app.DB}
		lastUpdate, err := commentRepo.GetLastUpdate(article.ID)
		if err != nil {
			app.ServerError(w, err)

			return
		}

		templateData := views.NewArticlePageData(article, lastUpdate)
		templateData.AppendTitle(article.Title)

		err = views.WriteTemplate(w, "article.gohtml", templateData)
		if err != nil {
			app.ServerError(w, err)
		}
	}
}