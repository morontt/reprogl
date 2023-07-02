package handlers

import (
	"net/http"

	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/views"
)

func LoginAction(app *container.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		templateData := views.NewLoginPageData()

		err := views.WriteTemplate(w, "login.gohtml", templateData)
		if err != nil {
			app.ServerError(w, err)

			return
		}
	}
}
