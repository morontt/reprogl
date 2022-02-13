package handlers

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/models"
	"xelbot.com/reprogl/models/repositories"
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

		articlesPaginator.URLGenerator = indexPaginationURLs(app.Router)

		tagRepo := repositories.TagRepository{DB: app.DB}
		err = tagRepo.PopulateTagsToArticles(articlesPaginator.Items)
		if err != nil {
			app.ServerError(w, err)

			return
		}

		templateData := views.NewIndexPageData(articlesPaginator)
		if page > 1 {
			templateData.AppendTitle(fmt.Sprintf("Страница %d", page))
		}

		doESI(w)
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
			url, err := app.Router.Get("category-first").URL("slug", slug)
			if err != nil {
				panic(err)
			}
			http.Redirect(w, r, url.String(), 301)

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

		articlesPaginator.URLGenerator = categoryPaginationURLs(app.Router, slug)

		tagRepo := repositories.TagRepository{DB: app.DB}
		err = tagRepo.PopulateTagsToArticles(articlesPaginator.Items)
		if err != nil {
			app.ServerError(w, err)

			return
		}

		templateData := views.NewCategoryPageData(articlesPaginator, category)
		browserTitle := fmt.Sprintf("Категория \"%s\"", category.Name)
		if page > 1 {
			browserTitle += fmt.Sprintf(". Страница %d", page)
		}
		templateData.AppendTitle(browserTitle)

		doESI(w)
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
			url, err := app.Router.Get("tag-first").URL("slug", slug)
			if err != nil {
				panic(err)
			}
			http.Redirect(w, r, url.String(), 301)

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

		articlesPaginator.URLGenerator = tagPaginationURLs(app.Router, slug)

		err = tagRepo.PopulateTagsToArticles(articlesPaginator.Items)
		if err != nil {
			app.ServerError(w, err)

			return
		}

		templateData := views.NewCategoryPageData(articlesPaginator, tag)
		browserTitle := fmt.Sprintf("Тег \"%s\"", tag.Name)
		if page > 1 {
			browserTitle += fmt.Sprintf(". Страница %d", page)
		}
		templateData.AppendTitle(browserTitle)

		doESI(w)
		err = views.WriteTemplate(w, "index.gohtml", templateData)
		if err != nil {
			app.ServerError(w, err)
		}
	}
}

func indexPaginationURLs(router *mux.Router) models.URLGenerator {
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

		url, err := router.Get("blog-page").URL("page", strconv.Itoa(page))
		if err != nil {
			panic(err)
		}

		return url.String()
	}
}

func categoryPaginationURLs(router *mux.Router, slug string) models.URLGenerator {
	return paginationURLsWithSlug(router, slug, "category-first", "category")
}

func tagPaginationURLs(router *mux.Router, slug string) models.URLGenerator {
	return paginationURLsWithSlug(router, slug, "tag-first", "tag")
}

func paginationURLsWithSlug(router *mux.Router, slug, firstRouteName, othersRouteName string) models.URLGenerator {
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

		url, err := router.Get(routeName).URL(pairs...)
		if err != nil {
			panic(err)
		}

		return url.String()
	}
}
