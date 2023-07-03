package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/xelbot/yetacache"
	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/session"
	"xelbot.com/reprogl/views"
)

func LoginAction(app *container.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var csrfToken string
		var found bool

		cache := app.GetStringCache()

		if cookie, errNoCookie := r.Cookie(session.CsrfCookie); errNoCookie != nil {
			csrfToken = generateCsrfPair(w, cache)
		} else {
			csrfTokenKey := cookie.Value
			if csrfToken, found = cache.Get(csrfTokenKey); !found {
				csrfToken = generateCsrfPair(w, cache)
			}
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

func generateCsrfPair(w http.ResponseWriter, cache *yetacache.Cache[string, string]) string {
	csrfToken := generateToken()
	csrfTokenKey := generateToken()

	cache.Set(csrfTokenKey, csrfToken, 30*time.Minute)
	session.WriteSessionCookie(w, session.CsrfCookie, csrfTokenKey, "/login")

	return csrfToken
}

func generateToken() string {
	nonce := make([]byte, 18)
	_, err := rand.Read(nonce)
	if err != nil {
		panic(err)
	}

	return base64.URLEncoding.EncodeToString(nonce)
}
