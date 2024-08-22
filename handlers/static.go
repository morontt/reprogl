package handlers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/models/repositories"
	"xelbot.com/reprogl/views"
)

func InfoAction(app *container.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		templateData := views.NewInfoPageData()
		templateData.SetCanonical(container.GenerateAbsoluteURL("info-page"))
		err := views.WriteTemplateWithContext(r.Context(), w, "about.gohtml", templateData)
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

		articleRepo := repositories.ArticleRepository{DB: app.DB}
		monthArticles, err := articleRepo.GetMostVisitedArticlesOfMonth()
		if err != nil {
			app.ServerError(w, err)

			return
		}

		allTimeArticles, err := articleRepo.GetMostVisitedArticles()
		if err != nil {
			app.ServerError(w, err)

			return
		}

		templateData := views.NewStatisticsPageData()
		templateData.Commentators = commentators
		templateData.MonthArticles = monthArticles
		templateData.AllTimeArticles = allTimeArticles
		templateData.SetCanonical(container.GenerateAbsoluteURL("statistics"))

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
			"User-agent: *\n\nSitemap: https://%s/sitemap.xml\n",
			cfg.Host)
	}

	cacheControl(w, container.RobotsTxtTTL)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(body))
}

func HumansTXTAction(w http.ResponseWriter, r *http.Request) {
	var lastUpdateStr string
	lastUpdate, err := time.Parse(time.RFC3339, container.BuildTime)
	if err != nil {
		lastUpdateStr = container.BuildTime
	} else {
		lastUpdateStr = lastUpdate.Format("2006/01/02")
	}

	cfg := container.GetConfig()

	body := fmt.Sprintf(
		`/* TEAM */
	Developer: %s
	Contact: %s
	Telegram: @morontt
	From: %s

/* THANKS */
	Inspirer: Alex Edwards
	Site: https://www.alexedwards.net/
	From: Austria

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
	Server: NGINX + Varnish + Golang custom application
	Powered by: Go %s`,
		cfg.Author.FullName,
		cfg.Author.Email,
		cfg.AuthorLocationEn,
		lastUpdateStr,
		container.GoVersionNumbers,
	)

	cacheControl(w, container.StatisticsTTL)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(body + "\n"))
}

func FavIconAction(w http.ResponseWriter, _ *http.Request) {
	body, err := os.ReadFile("./public/favicon.ico")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	cacheControl(w, container.RobotsTxtTTL)
	w.Header().Set("Content-Type", "image/x-icon")
	w.Write(body)
}
