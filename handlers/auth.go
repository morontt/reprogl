package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/xelbot/yetacache"
	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/models"
	"xelbot.com/reprogl/models/repositories"
	"xelbot.com/reprogl/services/auth"
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

		saveLoginReferer(w, r)
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

		user, err := auth.HandleLoginPassword(app, r.PostForm.Get("username"), r.PostForm.Get("password"))
		if err != nil {
			deleteCsrfCookie(w)
			if authError, found := err.(auth.Error); found {
				session.Put(r.Context(), session.FlashErrorKey, err.Error())
				app.InfoLog.Println("[AUTH] " + authError.InfoLogMessage())
				http.Redirect(w, r, container.GenerateURL("login"), http.StatusSeeOther)
			} else {
				app.ServerError(w, err)
			}

			return
		}

		session.Put(r.Context(), session.FlashSuccessKey, fmt.Sprintf("Привет, %s :)", user.Username))
		app.InfoLog.Printf("[AUTH] success for \"%s\"\n", user.Username)
		authSuccess(user, app, container.RealRemoteAddress(r), r.Context())

		var redirectUrl string
		if redirectUrl, found = popLoginReferer(w, r); !found {
			redirectUrl = "/"
		}

		deleteCsrfCookie(w)
		http.Redirect(w, r, redirectUrl, http.StatusFound)
	}
}

func LogoutAction(w http.ResponseWriter, r *http.Request) {
	session.Destroy(r.Context())
	http.Redirect(w, r, "/", http.StatusFound)
}

func AuthNavigation(app *container.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cacheControl(w, container.DefaultEsiTTL)
		templateData := views.NewAuthNavigationData()
		err := views.WriteTemplateWithContext(r.Context(), w, "auth-navigation.gohtml", templateData)
		if err != nil {
			app.ServerError(w, err)
		}
	}
}

func MenuAuth(app *container.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user *models.User
		if identity, ok := session.GetIdentity(r.Context()); ok {
			repo := repositories.UserRepository{DB: app.DB}
			user, _ = repo.Find(identity.ID)
		}

		templateData := views.NewMenuAuthData(user)
		cacheControl(w, container.DefaultEsiTTL)

		err := views.WriteTemplateWithContext(r.Context(), w, "menu-auth.gohtml", templateData)
		if err != nil {
			app.ServerError(w, err)
		}
	}
}

func generateCsrfPair(w http.ResponseWriter, cache *yetacache.Cache[string, string]) string {
	csrfToken := generateRandomToken()
	csrfTokenKey := generateRandomToken()

	cache.Set(csrfTokenKey, csrfToken, 30*time.Minute)
	session.WriteSessionCookie(w, session.CsrfCookie, csrfTokenKey, "/login")

	return csrfToken
}

func deleteCsrfCookie(w http.ResponseWriter) {
	session.DeleteCookie(w, session.CsrfCookie, "/login")
}
