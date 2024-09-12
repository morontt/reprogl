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
	"xelbot.com/reprogl/utils/hashid"
	"xelbot.com/reprogl/views"
)

func PageAction(app *container.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		slug := vars["slug"]

		var isAdmin bool
		if user, ok := session.GetIdentity(r.Context()); ok {
			isAdmin = user.IsAdmin()
		}

		repo := repositories.ArticleRepository{DB: app.DB}
		article, err := repo.GetBySlug(slug, isAdmin)
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

		templateData := views.NewArticlePageData(
			article,
			lastUpdate,
			r.Header.Get("Accept"),
		)
		cfg := container.GetConfig()
		templateData.AppendTitle(article.Title)
		templateData.SetCanonical(container.GenerateAbsoluteURL("article", "slug", slug))
		templateData.SetOpenGraphProperty("og:type", "article")
		templateData.SetOpenGraphProperty("article:published_time", article.CreatedAt.Format(time.RFC3339))
		templateData.SetOpenGraphProperty("article:modified_time", article.UpdatedAt.Format(time.RFC3339))
		templateData.SetOpenGraphProperty("article:author", cfg.Author.FullName)
		if article.HasImage() {
			image := article.SrcImageForOpenGraph()
			if image != nil {
				templateData.SetOpenGraphProperty("og:image", cfg.CDNBaseURL+"/uploads/"+image.Path)
				templateData.SetOpenGraphProperty("og:image:width", strconv.Itoa(image.Width))
				templateData.SetOpenGraphProperty("og:image:height", strconv.Itoa(image.Height))
			}
			if article.Alt.Valid && len(article.Alt.String) > 0 {
				templateData.SetOpenGraphProperty("og:image:alt", article.Alt.String)
			}
		}
		templateData.AuthorAvatar = models.AvatarLink(1, hashid.Male|hashid.User, 200)

		err = views.WriteTemplateWithContext(r.Context(), w, "article.gohtml", templateData)
		if err != nil {
			app.ServerError(w, err)
		}
	}
}
