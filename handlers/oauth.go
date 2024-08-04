package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/services/oauth"
	"xelbot.com/reprogl/session"
)

func OAuthLogin(app *container.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		providerName := vars["provider"]

		app.InfoLog.Println("[OAUTH] start authorization by " + providerName)
		oauthConfig, err := oauth.ConfigByProvider(providerName)
		if err != nil {
			app.NotFound(w)

			return
		}

		state := generateRandomToken()
		session.Put(r.Context(), session.OAuthStateKey, state)
		app.InfoLog.Println("[OAUTH] generate state: " + state)

		url := oauthConfig.AuthCodeURL(state)

		// http.Redirect(w, r, u, http.StatusTemporaryRedirect)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("URL: " + url))
	}
}
