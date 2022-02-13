package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
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
