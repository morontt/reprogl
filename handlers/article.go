package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/models"
	"xelbot.com/reprogl/models/repositories"
	"xelbot.com/reprogl/session"
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

		cache := app.GetIntCache()

		var recentID int
		var found bool
		if recentID, found = cache.Get("last_recent_id"); !found {
			recentID, err = repo.GetLastRecentPostsID()
			if err != nil {
				app.ServerError(w, err)

				return
			}
			cache.Set("last_recent_id", recentID, 24*time.Hour)
		}
		if article.ID >= recentID {
			article.RecentPostsID = strconv.Itoa(article.ID)
		} else {
			article.RecentPostsID = "0"
		}

		templateData := views.NewArticlePageData(article, lastUpdate, session.HasIdentity(r.Context()))
		templateData.AppendTitle(article.Title)

		err = views.WriteTemplateWithContext(r.Context(), w, "article.gohtml", templateData)
		if err != nil {
			app.ServerError(w, err)
		}
	}
}
