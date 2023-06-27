package handlers

import (
	"encoding/base64"
	"github.com/gorilla/mux"
	"net/http"

	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/utils/hashid"
)

func AvatarGenerator(app *container.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		hash := vars["hash"]

		_, err := hashid.Decode(hash)
		if err != nil {
			app.ClientError(w, http.StatusBadRequest)
			return
		}

		pixel := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNkYAAAAAYAAjCB0C8AAAAASUVORK5CYII="
		body, err := base64.StdEncoding.DecodeString(pixel)
		if err != nil {
			panic(err)
		}

		w.Header().Set("Content-Type", "image/png")
		w.Write(body)
	}
}
