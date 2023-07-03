package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"xelbot.com/reprogl/session"

	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/views"
)

func LoginAction(app *container.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var csrfToken string
		var found bool

		if csrfToken, found = session.GetString(r.Context(), "csrf_token"); !found {
			nonce := make([]byte, 18)
			_, err := rand.Read(nonce)
			if err != nil {
				panic(err)
			}

			csrfToken = base64.URLEncoding.EncodeToString(nonce)
			session.Put(r.Context(), "csrf_token", csrfToken)
		}

		templateData := views.NewLoginPageData(csrfToken)
		err := views.WriteTemplate(w, "login.gohtml", templateData)
		if err != nil {
			app.ServerError(w, err)

			return
		}
	}
}

func LoginLogoutLinks(app *container.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//templateData := &views.FragmentRecentPostsData{RecentPosts: articles}
		//cacheControl(w, container.DefaultEsiTTL)
		err := views.WriteTemplate(w, "login-logout.gohtml", nil)
		if err != nil {
			app.ServerError(w, err)
		}
	}
}
