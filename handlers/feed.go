package handlers

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"runtime"
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

		f := urlBySlugGenerator(app.Router)

		for _, location := range *articles {
			location.URL = f(location.Slug)
		}

		urlSet := models.SitemapURLSet{Items: articles}

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

func FeedAction(app *container.Application, feedType int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var feed models.FeedGeneratorInterface

		repo := repositories.ArticleRepository{DB: app.DB}
		articles, err := repo.GetFeedCollection()
		if err != nil {
			app.ServerError(w, err)

			return
		}

		f := urlBySlugGenerator(app.Router)

		for _, location := range *articles {
			location.URL = f(location.Slug)
		}

		switch feedType {
		case models.AtomFeedType:
			feed = models.CreateFeed(new(models.Atom), channelData(articles))
		case models.RssFeedType:
			feed = models.CreateFeed(new(models.Rss), channelData(articles))
		default:
			app.ServerError(w, errors.New("undefined feed type"))

			return
		}

		bytes, err := feed.AsXML()
		if err != nil {
			app.ServerError(w, err)

			return
		}

		w.Header().Set("Content-Type", feed.ContentType())
		w.Header().Set("Cache-Control", "max-age=5")
		_, err = w.Write([]byte(xml.Header))
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

func channelData(items *models.FeedItemList) models.FeedChannelData {
	cfg := container.GetConfig()

	channel := models.FeedChannelData{
		Title:       cfg.Host,
		Link:        "https://" + cfg.Host + "/",
		Description: cfg.Host + " - последние записи",
		Language:    "ru-ru",
		Charset:     "utf-8",
		Author:      cfg.Author,
		Email:       cfg.AdminEmail,
		Generator: fmt.Sprintf(
			"Reprogl/%s (%s)",
			container.GitRevision,
			runtime.Version()),
		FeedItems: items,
	}

	return channel
}
