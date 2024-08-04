package handlers

import (
	"context"
	"fmt"
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

		app.InfoLog.Println("[OAUTH] start authorization by: " + providerName)
		oauthConfig, err := oauth.ConfigByProvider(providerName)
		if err != nil {
			app.NotFound(w)

			return
		}

		state := generateRandomToken()
		session.Put(r.Context(), session.OAuthStateKey, state)

		http.Redirect(w, r, oauthConfig.AuthCodeURL(state), http.StatusTemporaryRedirect)
	}
}

func OAuthCallback(app *container.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		providerName := vars["provider"]

		app.InfoLog.Println("[OAUTH] callback from: " + providerName)
		oauthConfig, err := oauth.ConfigByProvider(providerName)
		if err != nil {
			app.NotFound(w)

			return
		}

		state, _ := session.Pop[string](r.Context(), session.OAuthStateKey)
		stateFromRequest := r.FormValue("state")

		if len(state) == 0 || len(stateFromRequest) == 0 || stateFromRequest != state {
			app.InfoLog.Println("[OAUTH] Invalid state")
			app.ClientError(w, http.StatusBadRequest)

			return
		}

		code := r.FormValue("code")
		if len(code) == 0 {
			app.InfoLog.Println("[OAUTH] Empty code")
			app.ClientError(w, http.StatusBadRequest)

			return
		}

		token, err := oauthConfig.Exchange(context.Background(), code)
		if err != nil {
			app.ServerError(w, err)

			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(fmt.Sprintf("%+v\n", token)))
	}
}
