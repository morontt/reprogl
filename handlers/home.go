package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/models"
	"xelbot.com/reprogl/models/repositories"
	"xelbot.com/reprogl/session"
	"xelbot.com/reprogl/views"
)

func IndexAction(app *container.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		page, needsRedirect := pageOrRedirect(vars)
		if needsRedirect {
			http.Redirect(w, r, "/", 301)

			return
		}

		repo := repositories.ArticleRepository{DB: app.DB}
		articlesPaginator, err := repo.GetCollection(page)
		if err != nil {
			if errors.Is(err, models.RecordNotFound) {
				app.NotFound(w)
			} else {
				app.ServerError(w, err)
			}

			return
		}

		articlesPaginator.URLGenerator = indexPaginationURLs()

		tagRepo := repositories.TagRepository{DB: app.DB}
		err = tagRepo.PopulateTagsToArticles(articlesPaginator.Items)
		if err != nil {
			app.ServerError(w, err)

			return
		}

		flashSuccessMessage, _ := session.Pop[string](r.Context(), session.FlashSuccessKey)
		templateData := views.NewIndexPageData(articlesPaginator, flashSuccessMessage)
		if page > 1 {
			templateData.AppendTitle(fmt.Sprintf("Страница %d", page))
		}

		err = views.WriteTemplate(w, "index.gohtml", templateData)
		if err != nil {
			app.ServerError(w, err)
		}
	}
}

func CategoryAction(app *container.Application) http.HandlerFunc {
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
			http.Redirect(w, r, container.GenerateURL("category-first", "slug", slug), http.StatusMovedPermanently)

			return
		}

		repo := repositories.ArticleRepository{DB: app.DB}
		articlesPaginator, err := repo.GetCollectionByCategory(category, page)
		if err != nil {
			if errors.Is(err, models.RecordNotFound) {
				app.NotFound(w)
			} else {
				app.ServerError(w, err)
			}

			return
		}

		articlesPaginator.URLGenerator = categoryPaginationURLs(slug)

		tagRepo := repositories.TagRepository{DB: app.DB}
		err = tagRepo.PopulateTagsToArticles(articlesPaginator.Items)
		if err != nil {
			app.ServerError(w, err)

			return
		}

		templateData := views.NewCategoryPageData(articlesPaginator, category)

		err = views.WriteTemplate(w, "index.gohtml", templateData)
		if err != nil {
			app.ServerError(w, err)
		}
	}
}

func TagAction(app *container.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		slug := vars["slug"]

		tagRepo := repositories.TagRepository{DB: app.DB}
		tag, err := tagRepo.GetBySlug(slug)
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
			http.Redirect(w, r, container.GenerateURL("tag-first", "slug", slug), http.StatusMovedPermanently)

			return
		}

		repo := repositories.ArticleRepository{DB: app.DB}
		articlesPaginator, err := repo.GetCollectionByTag(tag, page)
		if err != nil {
			if errors.Is(err, models.RecordNotFound) {
				app.NotFound(w)
			} else {
				app.ServerError(w, err)
			}

			return
		}

		articlesPaginator.URLGenerator = tagPaginationURLs(slug)

		err = tagRepo.PopulateTagsToArticles(articlesPaginator.Items)
		if err != nil {
			app.ServerError(w, err)

			return
		}

		templateData := views.NewCategoryPageData(articlesPaginator, tag)

		err = views.WriteTemplate(w, "index.gohtml", templateData)
		if err != nil {
			app.ServerError(w, err)
		}
	}
}

func indexPaginationURLs() models.URLGenerator {
	return func(page int, dir models.PaginationDirection) string {
		switch dir {
		case models.PaginationNext:
			page++
		case models.PaginationPrev:
			page--
		}

		if page == 1 {
			return "/"
		}

		return container.GenerateURL("blog-page", "page", strconv.Itoa(page))
	}
}

func categoryPaginationURLs(slug string) models.URLGenerator {
	return paginationURLsWithSlug(slug, "category-first", "category")
}

func tagPaginationURLs(slug string) models.URLGenerator {
	return paginationURLsWithSlug(slug, "tag-first", "tag")
}

func paginationURLsWithSlug(slug, firstRouteName, othersRouteName string) models.URLGenerator {
	return func(page int, dir models.PaginationDirection) string {
		var (
			routeName string
			pairs     []string
		)

		switch dir {
		case models.PaginationNext:
			page++
		case models.PaginationPrev:
			page--
		}

		pairs = append(pairs, "slug")
		pairs = append(pairs, slug)

		if page == 1 {
			routeName = firstRouteName
		} else {
			routeName = othersRouteName
			pairs = append(pairs, "page")
			pairs = append(pairs, strconv.Itoa(page))
		}

		return container.GenerateURL(routeName, pairs...)
	}
}
