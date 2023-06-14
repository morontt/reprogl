package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/models/repositories"
	"xelbot.com/reprogl/views"
)

func InfoAction(app *container.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := views.WriteTemplate(w, "info.gohtml", views.NewInfoPageData())
		if err != nil {
			app.ServerError(w, err)

			return
		}
	}
}

func StatisticsAction(app *container.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repo := repositories.CommentRepository{DB: app.DB}
		commentators, err := repo.GetMostActiveCommentators()
		if err != nil {
			app.ServerError(w, err)

			return
		}

		templateData := views.NewStatisticsPageData()
		templateData.Commentators = commentators

		cacheControl(w, container.StatisticsTTL)
		err = views.WriteTemplate(w, "statistics.gohtml", templateData)
		if err != nil {
			app.ServerError(w, err)

			return
		}
	}
}

func RobotsTXTAction(w http.ResponseWriter, r *http.Request) {
	var body string

	if container.IsCDN(r) {
		body = "User-agent: *\n\nDisallow: /\n"
	} else {
		cfg := container.GetConfig()
		body = fmt.Sprintf(
			"User-agent: *\n\nHost: %s\nSitemap: https://%s/sitemap.xml\n",
			cfg.Host,
			cfg.Host)
	}

	cacheControl(w, container.RobotsTxtTTL)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(body))
}

func FavIconAction(w http.ResponseWriter, _ *http.Request) {
	var icon string
	icon = `
		AAABAAEAEBAQAAEABAAoAQAAFgAAACgAAAAQAAAAIAAAAAEABAAAAAAAAAAAAAAAAAAAAAAAAAAA
		AAAAAAAAAAAA////AACAAAAAgIAAgAAAAIAAgACAgAAAgICAAOPj4wAAAP8AAP8AAAD//wD/AAAA
		/wD/AP//tQD///8AAqKgAAAAoqIAKgiIiIgKKgAAiACIiICiAAiIiACIiAoAAAiIiIiIAgAe4IiI
		iIgAABHgiAAIiIAAERCAHuCIgAAACIAR4IiAAAiIgBEQiIAACAiIAAiIAAAAgIiIiIgAAAgIiIiI
		iAAAAICAgIiAAAAAAICAiAAAAAAAAAAAAACAAAAAwAAAAOAAAADAAAAAwAAAAIABAACAAAAAgAAA
		AMAAAADAAAAAwAEAAMABAADAAQAA4AMAAPAHAAD4DwAA`

	body, _ := base64.StdEncoding.DecodeString(icon)

	cacheControl(w, container.RobotsTxtTTL)
	w.Header().Set("Content-Type", "image/x-icon")
	w.Write(body)
}
