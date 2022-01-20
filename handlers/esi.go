package handlers

import (
	"net/http"
	"xelbot.com/reprogl/config"
	"xelbot.com/reprogl/models/repositories"
	"xelbot.com/reprogl/views"
)

func CategoriesFragment(app *config.Application) http.HandlerFunc {
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
