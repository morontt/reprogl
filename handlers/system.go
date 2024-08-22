package handlers

import (
	"fmt"
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

func HeadersDebug(w http.ResponseWriter, r *http.Request) {
	var body string
	for name, values := range r.Header {
		for _, value := range values {
			body += fmt.Sprintf("%s: %s\n", name, value)
		}
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(body))
}
