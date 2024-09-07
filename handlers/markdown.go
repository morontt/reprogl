package handlers

import (
	"errors"
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/views"
)

func MarkdownAction(app *container.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		content, err := views.MarkdownToHTML(chi.URLParam(r, "filename"))
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				app.NotFound(w)
			} else {
				app.ServerError(w, err)
			}

			return
		}

		cacheControl(w, container.StaticFileTTL)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, err = w.Write(content)
		if err != nil {
			app.ServerError(w, err)

			return
		}
	}
}
