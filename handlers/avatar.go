package handlers

import (
	"errors"
	"github.com/gorilla/mux"
	"image/png"
	"net/http"

	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/models"
	"xelbot.com/reprogl/utils/avatar"
	"xelbot.com/reprogl/utils/hashid"
)

func AvatarGenerator(app *container.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		hash := vars["hash"]

		hashModel, err := hashid.Decode(hash)
		if err != nil {
			app.ClientError(w, http.StatusBadRequest)
			return
		}

		img, err := avatar.GenerateAvatar(hashModel, app)
		if err != nil {
			if errors.Is(err, models.RecordNotFound) {
				app.NotFound(w)
			} else {
				app.ServerError(w, err)
			}

			return
		}

		cacheControl(w, container.AvatarTTL)
		w.Header().Set("Content-Type", "image/png")
		png.Encode(w, img)
	}
}
