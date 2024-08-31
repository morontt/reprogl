package handlers

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/models"
	"xelbot.com/reprogl/models/repositories"
	"xelbot.com/reprogl/session"
	"xelbot.com/reprogl/utils/hashid"
	"xelbot.com/reprogl/views"
)

const successUnsubscribe = "success.unsubscribe"

func EmailUnsubscribe(app *container.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		hash := vars["hash"]

		data, err := hashid.Decode(hash, false)
		if err != nil {
			app.NotFound(w)

			return
		}

		repo := repositories.EmailSubscriptionRepository{DB: app.DB}
		settings, err := repo.Find(data.ID)
		if err != nil {
			if errors.Is(err, models.RecordNotFound) {
				app.NotFound(w)
			} else {
				app.ServerError(w, err)
			}

			return
		}

		_, ok := session.Pop[string](r.Context(), successUnsubscribe)

		templateData := views.NewUnsubscribePageData(
			settings,
			models.AvatarLink(3, hashid.Male|hashid.User, 160),
			ok || settings.BlockSending,
		)
		err = views.WriteTemplate(w, "unsubscribe.gohtml", templateData)
		if err != nil {
			app.ServerError(w, err)

			return
		}
	}
}

func EmailUnsubscribePost(app *container.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		hash := vars["hash"]

		data, err := hashid.Decode(hash, false)
		if err != nil {
			app.NotFound(w)

			return
		}

		repo := repositories.EmailSubscriptionRepository{DB: app.DB}
		err = repo.Unsubscribe(data.ID)
		if err != nil {
			app.ServerError(w, err)

			return
		}

		session.Put(r.Context(), successUnsubscribe, "*")

		http.Redirect(w, r, r.URL.Path, http.StatusFound)
	}
}
