package middlewares

import (
	"net/http"

	"xelbot.com/reprogl/container"
	pkghttp "xelbot.com/reprogl/http"
	"xelbot.com/reprogl/utils/tracking"
)

func Track(next http.Handler, app *container.Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		activity := tracking.CreateActivity(r)
		next.ServeHTTP(w, r)

		if activity != nil {
			lrw, ok := w.(pkghttp.LogResponseWriter)
			if ok {
				activity.Status = lrw.Status()
			}

			go tracking.SaveActivity(activity, app)
		}
	})
}
