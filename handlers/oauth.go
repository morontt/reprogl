package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"xelbot.com/reprogl/api/backend"
	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/models/repositories"
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

		saveLoginReferer(w, r.Referer())

		state := generateRandomToken()
		session.Put(r.Context(), session.OAuthStateKey, state)

		verifier := oauth2.GenerateVerifier()
		session.Put(r.Context(), session.OAuthVerifierKey, verifier)

		http.Redirect(w, r, oauthConfig.AuthCodeURL(state, oauth2.S256ChallengeOption(verifier)), http.StatusFound)
	}
}

func OAuthCallback(app *container.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		providerName := vars["provider"]

		app.InfoLog.Println("[OAUTH] callback from: " + providerName)
		if !oauth.SupportedProvider(providerName) {
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

		var found bool
		verifier, found := session.Pop[string](r.Context(), session.OAuthVerifierKey)
		if !found {
			app.ServerError(w, errors.New("[OAUTH] PKCE verifier not found"))

			return
		}

		code := r.FormValue("code")
		if len(code) == 0 {
			errorCode := r.FormValue("error")
			errorDescription := r.FormValue("error_description")
			if len(errorCode) > 0 {
				app.InfoLog.Printf("[OAUTH] Error code: %s, description: %s\n", errorCode, errorDescription)
			} else {
				app.InfoLog.Println("[OAUTH] Error: empty code")
			}
			app.ClientError(w, http.StatusBadRequest)

			return
		}

		userData, err := oauth.UserDataByCode(providerName, code, verifier)
		if err != nil {
			app.ServerError(w, err)

			return
		}

		userDataDTO := backend.ExternalUserDTO{
			UserData:  userData,
			UserAgent: r.UserAgent(),
			IP:        container.RealRemoteAddress(r),
		}

		apiResponse, err := backend.SendUserData(userDataDTO)
		if err != nil {
			app.ServerError(w, err)

			return
		}

		if apiResponse.Violations != nil && len(apiResponse.Violations) > 0 {
			errorMessage := "[OAUTH] user validation error:\n"
			for _, formError := range apiResponse.Violations {
				app.InfoLog.Printf("[OAUTH] validation error: %s - %s\n", formError.Path, formError.Message)
				errorMessage += fmt.Sprintf("%s: %s\n", formError.Path, formError.Message)
			}

			app.ServerError(w, errors.New(errorMessage))

			return
		}

		if apiResponse.User != nil {
			session.Put(r.Context(), session.FlashSuccessKey, fmt.Sprintf("Привет, %s :)", apiResponse.User.Nickname()))

			repo := repositories.UserRepository{DB: app.DB}
			user, err := repo.GetLoggedUserByUsername(apiResponse.User.Username)
			if err != nil {
				app.ServerError(w, err)

				return
			}

			app.InfoLog.Printf("[OAUTH] success for \"%s\"\n", user.Username)
			authSuccess(user, app, container.RealRemoteAddress(r), r.Context())
		}

		var redirectUrl string
		if redirectUrl, found = popLoginReferer(w, r); !found {
			redirectUrl = "/"
		}

		http.Redirect(w, r, redirectUrl, http.StatusFound)
	}
}
