package handlers

import (
	"errors"
	"net/http"

	"xelbot.com/reprogl/api/backend"
	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/models"
	"xelbot.com/reprogl/models/repositories"
	"xelbot.com/reprogl/session"
	"xelbot.com/reprogl/views"
)

type updateProfileResponse struct {
	Valid  bool                `json:"valid"`
	Errors []backend.FormError `json:"errors,omitempty"`
}

func ProfileAction(app *container.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO move to middleware
		if !session.HasIdentity(r.Context()) {
			http.Redirect(w, r, container.GenerateURL("login"), http.StatusFound)
			return
		}

		var user *models.User
		if identity, ok := session.GetIdentity(r.Context()); ok {
			repo := repositories.UserRepository{DB: app.DB}
			user, _ = repo.Find(identity.ID)
		}

		if user == nil {
			app.ServerError(w, errors.New("profile logic error: user is null"))
			return
		}

		subscriptionSettingsRepo := repositories.EmailSubscriptionRepository{DB: app.DB}
		subscrSettings, err := subscriptionSettingsRepo.FindOrCreate(user.Email, models.SubscriptionReplyComment)
		if err != nil {
			app.ServerError(w, err)
			return
		}

		templateData := views.NewProfilePageData(user, subscrSettings)
		err, wh := views.WriteTemplate(w, "profile.gohtml", templateData)
		if err != nil {
			if wh {
				app.LogError(err)
			} else {
				app.ServerError(w, err)
			}
		}
	}
}

func UpdateProfile(app *container.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO move to middleware
		if !session.HasIdentity(r.Context()) {
			http.Redirect(w, r, container.GenerateURL("login"), http.StatusFound)
			return
		}

		var user *models.User
		if identity, ok := session.GetIdentity(r.Context()); ok {
			repo := repositories.UserRepository{DB: app.DB}
			user, _ = repo.Find(identity.ID)
		}

		if user == nil {
			app.ServerError(w, errors.New("profile logic error: user is null"))
			return
		}

		emailFromForm := r.PostFormValue("email")
		if len(emailFromForm) == 0 && !user.HasEmail() {
			emailFromForm = user.Email
		}

		profileData := backend.ProfileDTO{
			ID:          user.ID,
			Role:        user.Role,
			Username:    r.PostFormValue("username"),
			Email:       emailFromForm,
			DisplayName: r.PostFormValue("displayName"),
		}

		if r.PostFormValue("gender") == "male" {
			profileData.IsMale = true
		} else {
			profileData.IsMale = false
		}

		var responseData any
		statusCode := http.StatusOK

		apiResponse, err := backend.SendProfileData(profileData)
		if err != nil {
			app.ServerError(w, err)

			return
		}

		if apiResponse != nil {
			result := updateProfileResponse{
				Valid:  true,
				Errors: apiResponse.Violations,
			}

			if len(apiResponse.Violations) > 0 {
				result.Valid = false
			}

			responseData = result

			if apiResponse.User != nil {
				subscriptionSettingsRepo := repositories.EmailSubscriptionRepository{DB: app.DB}
				subscrSettings, err := subscriptionSettingsRepo.FindOrCreate(
					apiResponse.User.Email,
					models.SubscriptionReplyComment)
				if err != nil {
					app.ServerError(w, err)
					return
				}

				formValue := r.PostFormValue("reply_subscribe")
				if formValue != "" {
					if subscrSettings.BlockSending {
						err = subscriptionSettingsRepo.Subscribe(subscrSettings.ID)
					}
				} else {
					if !subscrSettings.BlockSending {
						err = subscriptionSettingsRepo.Unsubscribe(subscrSettings.ID)
					}
				}
				if err != nil {
					app.ServerError(w, err)
					return
				}
			}
		}

		jsonResponse(w, statusCode, responseData)
	}
}
