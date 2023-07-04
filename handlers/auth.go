package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"time"

	"github.com/xelbot/yetacache"
	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/models"
	"xelbot.com/reprogl/models/repositories"
	"xelbot.com/reprogl/security"
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

func LoginCheck(app *container.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var csrfToken string
		var found bool

		if session.HasIdentity(r.Context()) {
			session.ClearIdentity(r.Context())
		}

		cache := app.GetStringCache()
		if cookie, errNoCookie := r.Cookie(session.CsrfCookie); errNoCookie == nil {
			csrfTokenKey := cookie.Value
			if csrfToken, found = cache.Get(csrfTokenKey); !found {
				deleteCsrfCookie(w)
				http.Redirect(w, r, container.GenerateURL("login"), http.StatusSeeOther)
				return
			}
		}

		err := r.ParseForm()
		if err != nil {
			deleteCsrfCookie(w)
			app.ClientError(w, http.StatusBadRequest)
			return
		}

		formCsrfToken := r.PostForm.Get("_csrf_token")
		if formCsrfToken != csrfToken {
			deleteCsrfCookie(w)
			app.InfoLog.Println("[AUTH] wrong CSRF-token")
			http.Redirect(w, r, container.GenerateURL("login"), http.StatusSeeOther)
			return
		}

		username := r.PostForm.Get("username")
		password := r.PostForm.Get("password")

		repo := repositories.UserRepository{DB: app.DB}
		user, err := repo.GetLoggedUserByUsername(username)
		if err != nil {
			deleteCsrfCookie(w)
			if errors.Is(err, models.RecordNotFound) {
				// TODO flash message
				app.InfoLog.Printf("[AUTH] user \"%s\" not found\n", username)
				http.Redirect(w, r, container.GenerateURL("login"), http.StatusSeeOther)
			} else {
				app.ServerError(w, err)
			}

			return
		}

		passwordHash := security.EncodePassword(password, user.Salt)
		if passwordHash != user.PasswordHash {
			// TODO flash message, random pause
			app.InfoLog.Printf("[AUTH] invalid password for \"%s\"\n", username)
		} else {
			app.InfoLog.Printf("[AUTH] success for \"%s\"\n", username)
			authSuccess(user, app, r.Context())
		}

		deleteCsrfCookie(w)
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func LoginLogoutLinks(app *container.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := views.WriteTemplate(w, "login-logout.gohtml", nil)
		if err != nil {
			app.ServerError(w, err)
		}
	}
}

func authSuccess(user *models.LoggedUser, app *container.Application, ctx context.Context) {
	session.SetIdentity(ctx, security.CreateIdentity(user))
}

func generateCsrfPair(w http.ResponseWriter, cache *yetacache.Cache[string, string]) string {
	csrfToken := generateToken()
	csrfTokenKey := generateToken()

	cache.Set(csrfTokenKey, csrfToken, 30*time.Minute)
	session.WriteSessionCookie(w, session.CsrfCookie, csrfTokenKey, "/login")

	return csrfToken
}

func deleteCsrfCookie(w http.ResponseWriter) {
	session.DeleteCookie(w, session.CsrfCookie, "/login")
}

func generateToken() string {
	nonce := make([]byte, 18)
	_, err := rand.Read(nonce)
	if err != nil {
		panic(err)
	}

	return base64.URLEncoding.EncodeToString(nonce)
}
