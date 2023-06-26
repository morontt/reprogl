package handlers

import (
	"net/http"

	"xelbot.com/reprogl/container"
)

func PurgeCache(app *container.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app.GetIntCache().Clear()
		app.InfoLog.Println("[CACHE] integer cache was cleared")

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Cache was cleared\n"))
	}
}
