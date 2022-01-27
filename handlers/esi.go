package handlers

import (
	"net/http"
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
		}

		templateData := &views.FragmentCategoriesData{Categories: categories}

		err = views.RenderTemplate(w, "categories.gohtml", templateData)
		if err != nil {
			app.ServerError(w, err)
		}
	}
}

func CommentsFragment(_ *container.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
	}
}
