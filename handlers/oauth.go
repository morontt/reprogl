package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"golang.org/x/oauth2"
	"xelbot.com/reprogl/api/backend"
	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/models/repositories"
	"xelbot.com/reprogl/services/oauth"
	"xelbot.com/reprogl/session"
	"xelbot.com/reprogl/views"
)

type oauthCallbackState struct {
	Status   string `json:"status"`
	UserName string `json:"username,omitempty"`
	NickName string `json:"nickname,omitempty"`
}

type oauthStateResponse struct {
	Status      string `json:"status"`
	RedirectURL string `json:"redirect_url,omitempty"`
}

func OAuthLogin(app *container.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		providerName := chi.URLParam(r, "provider")

		app.InfoLog.Println("[OAUTH] start authorization by: " + providerName)
		oauthConfig, err := oauth.ConfigByProvider(providerName)
		if err != nil {
			app.NotFound(w)

			return
		}

		saveLoginReferer(w, r)

		state := generateRandomToken()
		session.Put(r.Context(), session.OAuthStateKey, state)

		verifier := oauth2.GenerateVerifier()
		session.Put(r.Context(), session.OAuthVerifierKey, verifier)

		redirectURL := oauthConfig.AuthCodeURL(state, oauth2.S256ChallengeOption(verifier))
		app.InfoLog.Println("[OAUTH] redirect to: " + redirectURL)

		http.Redirect(w, r, redirectURL, http.StatusFound)
	}
}

func OAuthCallback(app *container.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		providerName := chi.URLParam(r, "provider")

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

		additional := make(map[string]string)
		for _, key := range oauth.AdditionalParams(providerName) {
			additional[key] = r.FormValue(key)
		}

		requestID := generateRandomToken()
		go asyncCallback(requestID, providerName, code, verifier, r.UserAgent(), container.RealRemoteAddress(r), additional, app)

		templateData := views.NewOauthPendingPageData(requestID)
		err := views.WriteTemplate(w, "oauth-pending.gohtml", templateData)
		if err != nil {
			app.ServerError(w, err)

			return
		}
	}
}

func asyncCallback(
	requestID,
	providerName,
	code,
	verifier,
	userAgent,
	ip string,
	additional map[string]string,
	app *container.Application,
) {
	cache := app.GetStringCache()
	cache.Set(requestID, `{"status":"pending"}`, time.Minute)

	userData, err := oauth.UserDataByCode(providerName, code, verifier, additional)
	if err != nil {
		oauthCallbackError(app, requestID, err)

		return
	}

	userDataDTO := backend.ExternalUserDTO{
		UserData:  userData,
		UserAgent: userAgent,
		IP:        ip,
	}

	apiResponse, err := backend.SendUserData(userDataDTO)
	if err != nil {
		oauthCallbackError(app, requestID, err)

		return
	}

	if len(apiResponse.Violations) > 0 {
		errorMessage := "[OAUTH] user validation error:\n"
		for _, formError := range apiResponse.Violations {
			app.InfoLog.Printf("[OAUTH] validation error: %s - %s\n", formError.Path, formError.Message)
			errorMessage += fmt.Sprintf("%s: %s\n", formError.Path, formError.Message)
		}

		oauthCallbackError(app, requestID, err)

		return
	}

	if apiResponse.User != nil {
		oauthState := oauthCallbackState{
			Status:   "ok",
			UserName: apiResponse.User.Username,
			NickName: apiResponse.User.Nickname(),
		}

		jsonBody, err := json.Marshal(oauthState)
		if err != nil {
			oauthCallbackError(app, requestID, err)

			return
		}
		cache.Set(requestID, string(jsonBody), time.Minute)
	}
}

func OAuthCheckState(app *container.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := chi.URLParam(r, "request_id")

		var stateString string
		var found bool

		cache := app.GetStringCache()
		if stateString, found = cache.Get(requestID); !found {
			app.InfoLog.Println("[OAUTH] requestID not found: " + requestID)
			app.NotFound(w)

			return
		}

		buf := []byte(stateString)
		if !json.Valid(buf) {
			app.ServerError(w, errors.New("[OAUTH] invalid JSON state"))

			return
		}

		var oauthState oauthCallbackState
		err := json.Unmarshal(buf, &oauthState)
		if err != nil {
			app.ServerError(w, err)

			return
		}

		responseData := oauthStateResponse{
			Status: oauthState.Status,
		}

		if oauthState.Status == "ok" && len(oauthState.UserName) > 0 {
			session.Put(r.Context(), session.FlashSuccessKey, fmt.Sprintf("Привет, %s :)", oauthState.NickName))

			repo := repositories.UserRepository{DB: app.DB}
			user, err := repo.GetLoggedUserByUsername(oauthState.UserName)
			if err != nil {
				app.ServerError(w, err)

				return
			}

			app.InfoLog.Printf("[OAUTH] success for \"%s\"\n", user.Username)
			authSuccess(user, app, container.RealRemoteAddress(r), r.Context())

			var redirectUrl string
			if redirectUrl, found = popLoginReferer(w, r); !found {
				redirectUrl = "/"
			}

			responseData.RedirectURL = redirectUrl
		}

		jsonResponse(w, http.StatusOK, responseData)
	}
}

func oauthCallbackError(app *container.Application, requestID string, err error) {
	app.LogError(err)

	cache := app.GetStringCache()
	cache.Set(requestID, `{"status":"error"}`, time.Minute)
}
