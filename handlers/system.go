package handlers

import (
	"fmt"
	"net/http"
	"sort"

	"xelbot.com/reprogl/container"
)

func PurgeCache(app *container.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app.GetIntCache().Clear()
		app.InfoLog.Println("[CACHE] integer cache was cleared")

		app.GetStringCache().Clear()
		app.InfoLog.Println("[CACHE] string cache was cleared")

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Cache was cleared\n"))
	}
}

func HeadersDebug(w http.ResponseWriter, r *http.Request) {
	keys := make([]string, 0, len(r.Header))

	for k := range r.Header {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var body string
	for _, key := range keys {
		for _, value := range r.Header[key] {
			body += fmt.Sprintf("%s: %s\n", key, value)
		}
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(body))
}
