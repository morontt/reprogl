package handlers

import (
	"net/http"
	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/models/repositories"
	"xelbot.com/reprogl/views"
)

func SitemapAction(app *container.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repo := repositories.ArticleRepository{DB: app.DB}
		articles, err := repo.GetSitemapCollection()
		if err != nil {
			app.ServerError(w, err)

			return
		}

		xml, err := views.RenderTemplate("sitemap.gohtml", &views.SitemapData{Articles: articles})
		if err != nil {
			app.ServerError(w, err)

			return
		}

		w.Header().Set("Content-Type", "application/xml")
		w.Header().Set("Cache-Control", "public, max-age=1")
		_, err = w.Write([]byte(xml))
		if err != nil {
			app.ServerError(w, err)
		}
	}
}
