package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

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

func HumansTXTAction(w http.ResponseWriter, r *http.Request) {
	var lastUpdateStr string
	lastUpdate, err := time.Parse(time.RFC1123, container.BuildTime)
	if err != nil {
		lastUpdateStr = container.BuildTime
	} else {
		lastUpdateStr = lastUpdate.Format("2006/01/02")
	}

	body := fmt.Sprintf(
		`/* TEAM */
	Developer: Alexander Harchenko
	Contact: morontt [at] gmail [dot] com
	Twitter: @morontt
	Telegram: @morontt
	From: Kharkov, Ukraine

/* THANKS */
	pxThemes: Anima Multipurpose Ghost Theme
	Site: https://themeforest.net/user/pxthemes
	Location: Poland

	Colorlib: Free 404 Error Page Templates
	Site: https://colorlib.com/wp/free-404-error-page-templates/
	Location: Latvia

	8biticon: Pixel Character Maker
	Site: https://8biticon.com/

/* SITE */
	Last update: %s
	Language: Russian
	Doctype: HTML5
	IDE: GoLand, VS Code, nano
	Server: NGINX + Varnish + Golang http.Server
	Powered by: Go %s`,
		lastUpdateStr,
		container.GoVersionNumbers,
	)

	cacheControl(w, container.StatisticsTTL)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(body + "\n"))
}

func FavIconAction(w http.ResponseWriter, _ *http.Request) {
	var icon string
	icon = "AAABAAEAEBAQAAEABAAoAQAAFgAAACgAAAAQAAAAIAAAAAEABAAAAAAAAAAAAAAAAAAAAAAAAAAA" +
		"AAAAAAAAAAAA////AACAAAAAgIAAgAAAAIAAgACAgAAAgICAAOPj4wAAAP8AAP8AAAD//wD/AAAA" +
		"/wD/AP//tQD///8AAqKgAAAAoqIAKgiIiIgKKgAAiACIiICiAAiIiACIiAoAAAiIiIiIAgAe4IiI" +
		"iIgAABHgiAAIiIAAERCAHuCIgAAACIAR4IiAAAiIgBEQiIAACAiIAAiIAAAAgIiIiIgAAAgIiIiI" +
		"iAAAAICAgIiAAAAAAICAiAAAAAAAAAAAAACAAAAAwAAAAOAAAADAAAAAwAAAAIABAACAAAAAgAAA" +
		"AMAAAADAAAAAwAEAAMABAADAAQAA4AMAAPAHAAD4DwAA"

	body, err := base64.StdEncoding.DecodeString(icon)
	if err != nil {
		panic(err)
	}

	cacheControl(w, container.RobotsTxtTTL)
	w.Header().Set("Content-Type", "image/x-icon")
	w.Write(body)
}
