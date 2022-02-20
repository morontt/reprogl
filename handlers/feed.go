package handlers

import (
	"encoding/xml"
	"net/http"
	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/models"
	"xelbot.com/reprogl/models/repositories"
)

func SitemapAction(app *container.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repo := repositories.ArticleRepository{DB: app.DB}
		articles, err := repo.GetSitemapCollection()
		if err != nil {
			app.ServerError(w, err)

			return
		}

		router := app.Router
		cfg := container.GetConfig()
		f := func(slug string) string {
			url, _ := router.Get("article").URL("slug", slug)

			return "https://" + cfg.Host + url.String()
		}

		for _, location := range *articles {
			location.URL = f(location.Slug)
		}

		urlSet := models.SitemapURLSet{
			Items:     articles,
			Namespace: "http://www.sitemaps.org/schemas/sitemap/0.9",
		}

		bytes, err := xml.Marshal(urlSet)
		if err != nil {
			app.ServerError(w, err)

			return
		}

		w.Header().Set("Content-Type", "application/xml")
		w.Header().Set("Cache-Control", "max-age=7200")
		_, err = w.Write([]byte(xml.Header + `<?xml-stylesheet type="text/xsl" href="/sitemap.xsl"?>` + "\n"))
		if err != nil {
			app.ServerError(w, err)

			return
		}

		_, err = w.Write(bytes)
		if err != nil {
			app.ServerError(w, err)
		}
	}
}
