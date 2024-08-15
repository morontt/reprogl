package handlers

import (
	"errors"
	"net/http"

	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/models"
	"xelbot.com/reprogl/models/repositories"
	"xelbot.com/reprogl/session"
	"xelbot.com/reprogl/views"
)

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

		templateData := views.NewProfilePageData(user)
		err := views.WriteTemplate(w, "profile.gohtml", templateData)
		if err != nil {
			app.ServerError(w, err)

			return
		}
	}
}
