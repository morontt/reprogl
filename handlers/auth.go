package handlers

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
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

		errorMessage, hasError := session.Pop[string](r.Context(), session.FlashErrorKey)

		templateData := views.NewLoginPageData(csrfToken, errorMessage, hasError)
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
				session.Put(r.Context(), session.FlashErrorKey, "Непонятная ошибка")
				app.InfoLog.Println("[AUTH] not found CSRF-token in cache")
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
			session.Put(r.Context(), session.FlashErrorKey, "Непонятная ошибка")
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
				session.Put(r.Context(), session.FlashErrorKey, "Недействительные логин/пароль")
				app.InfoLog.Printf("[AUTH] user \"%s\" not found\n", username)
				http.Redirect(w, r, container.GenerateURL("login"), http.StatusSeeOther)
			} else {
				app.ServerError(w, err)
			}

			return
		}

		passwordHash := security.EncodePassword(password, user.Salt)
		if subtle.ConstantTimeCompare([]byte(passwordHash), []byte(user.PasswordHash)) == 0 {
			session.Put(r.Context(), session.FlashErrorKey, "Недействительные логин/пароль")
			app.InfoLog.Printf("[AUTH] invalid password for \"%s\"\n", username)
		} else {
			session.Put(r.Context(), session.FlashSuccessKey, fmt.Sprintf("Привет, %s :)", username))
			app.InfoLog.Printf("[AUTH] success for \"%s\"\n", username)
			authSuccess(user, app, container.RealRemoteAddress(r), r.Context())
		}

		deleteCsrfCookie(w)
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func LogoutAction(w http.ResponseWriter, r *http.Request) {
	session.Destroy(r.Context())
	http.Redirect(w, r, "/", http.StatusFound)
}

func LoginLogoutLinks(app *container.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cacheControl(w, container.DefaultEsiTTL)
		templateData := views.NewAuthNavigationData(session.HasIdentity(r.Context()))
		err := views.WriteTemplate(w, "auth-navigation.gohtml", templateData)
		if err != nil {
			app.ServerError(w, err)
		}
	}
}

func authSuccess(user *models.LoggedUser, app *container.Application, ip string, ctx context.Context) {
	session.SetIdentity(ctx, security.CreateIdentity(user))

	repo := repositories.UserRepository{DB: app.DB}
	if err := repo.SaveLoginEvent(user.ID, ip); err != nil {
		app.LogError(err)
	}
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
