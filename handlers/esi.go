package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/models/repositories"
	"xelbot.com/reprogl/views"
)

func CategoriesFragment(app *container.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		categoryRepo := repositories.CategoryRepository{DB: app.DB}
		categories, err := categoryRepo.GetCategoryTree()
		if err != nil {
			app.ServerError(w, err)

			return
		}

		templateData := &views.FragmentCategoriesData{Categories: categories}

		cacheControl(w, container.DefaultEsiTTL)
		err = views.WriteTemplate(w, "categories.gohtml", templateData)
		if err != nil {
			app.ServerError(w, err)
		}
	}
}

func CommentsFragment(app *container.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		articleId, err := strconv.Atoi(vars["article_id"])
		if err != nil {
			app.ServerError(w, err)

			return
		}

		repo := repositories.CommentRepository{DB: app.DB}
		comments, err := repo.GetCollectionByArticleId(articleId)
		if err != nil {
			app.ServerError(w, err)

			return
		}

		templateData := &views.FragmentCommentsData{Comments: comments}

		cacheControl(w, container.DefaultEsiTTL)
		err = views.WriteTemplate(w, "comments.gohtml", templateData)
		if err != nil {
			app.ServerError(w, err)
		}
	}
}

func RecentPostsFragment(app *container.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		articleId, err := strconv.Atoi(vars["article_id"])
		if err != nil {
			app.ServerError(w, err)

			return
		}

		repo := repositories.ArticleRepository{DB: app.DB}
		articles, err := repo.GetRecentPostsCollection(articleId)
		if err != nil {
			app.ServerError(w, err)

			return
		}

		templateData := &views.FragmentRecentPostsData{RecentPosts: articles}

		cacheControl(w, container.DefaultEsiTTL)
		err = views.WriteTemplate(w, "recent-posts.gohtml", templateData)
		if err != nil {
			app.ServerError(w, err)
		}
	}
}

func DaysOfWarCounter(w http.ResponseWriter, r *http.Request) {
	start, _ := time.ParseInLocation("2006-01-02", "2022-02-24", time.Local)
	now := time.Now()

	days := 1 + (now.Unix()-start.Unix())/(24*3600)
	response := fmt.Sprintf("<!-- War in Ukraine: Thu, 24 Feb 2022 05:00:00 +0200, %d days of war -->", days)

	w.Header().Set("Cache-Control", "public")

	expires, _ := time.ParseInLocation(time.DateTime, now.Format(time.DateOnly)+" 23:59:59", time.Local)
	setExpires(w, expires)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(response))
}
