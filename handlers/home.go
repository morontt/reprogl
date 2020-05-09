package handlers

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"xelbot.com/reprogl/config"
	"xelbot.com/reprogl/models"
	"xelbot.com/reprogl/models/repositories"
	"xelbot.com/reprogl/views"
)

func IndexAction(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		page, needsRedirect := pageOrRedirect(vars)
		if needsRedirect {
			http.Redirect(w, r, "/", 301)

			return
		}

		repo := repositories.ArticleRepository{DB: app.DB}
		articles, err := repo.GetCollection(page)
		if err != nil {
			app.ServerError(w, err)

			return
		}

		templateData := views.NewIndexPageData(articles, page)
		if page > 1 {
			templateData.AppendTitle(fmt.Sprintf("Страница %d", page))
		}

		err = views.RenderTemplate(w, "index.gohtml", templateData)
		if err != nil {
			app.ServerError(w, err)
		}
	}
}

func CategoryAction(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		slug := vars["slug"]

		categoryRepo := repositories.CategoryRepository{DB: app.DB}
		category, err := categoryRepo.GetBySlug(slug)
		if err != nil {
			if errors.Is(err, models.RecordNotFound) {
				app.NotFound(w)
			} else {
				app.ServerError(w, err)
			}

			return
		}

		page, needsRedirect := pageOrRedirect(vars)
		if needsRedirect {
			url, err := app.Router.Get("category-first").URL("slug", slug)
			if err != nil {
				panic(err)
			}
			http.Redirect(w, r, url.String(), 301)

			return
		}

		repo := repositories.ArticleRepository{DB: app.DB}
		articles, err := repo.GetCollectionByCategory(category, page)
		if err != nil {
			app.ServerError(w, err)

			return
		}

		templateData := views.NewCategoryPageData(articles, category, page)
		browserTitle := fmt.Sprintf("Категория \"%s\"", category.Name)
		if page > 1 {
			browserTitle += fmt.Sprintf(". Страница %d", page)
		}
		templateData.AppendTitle(browserTitle)

		err = views.RenderTemplate(w, "index.gohtml", templateData)
		if err != nil {
			app.ServerError(w, err)
		}
	}
}

func TagAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	page, needsRedirect := pageOrRedirect(vars)
	if needsRedirect {
		http.Redirect(w, r, "/", 301)

		return
	}

	fmt.Fprintf(w, "Articles by tag, page %d", page)
}
