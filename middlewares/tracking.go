package middlewares

import (
	"net/http"
	"strconv"
	"time"

	"xelbot.com/reprogl/container"
	pkghttp "xelbot.com/reprogl/http"
	"xelbot.com/reprogl/services/ratelimit"
	"xelbot.com/reprogl/utils/tracking"
)

func Track(next http.Handler, app *container.Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var activityHash string

		activity := tracking.CreateActivity(r)
		if activity != nil {
			activityHash = ratelimit.ActivityHash(activity)
			if ratelimit.IsBlocked(activityHash) {
				app.InfoLog.Println("[TRACKING] block by rate-limit: " + activityHash)

				w.Header().Set("Content-Type", "text/plain")
				w.Header().Set("Retry-After", strconv.Itoa(int(ratelimit.BlockingTime/time.Second)))
				w.WriteHeader(http.StatusTooManyRequests)

				w.Write([]byte("Too many request"))

				return
			}
		}

		next.ServeHTTP(w, r)

		if activity != nil {
			lrw, ok := w.(pkghttp.LogResponseWriter)
			if ok {
				activity.Status = lrw.Status()
				activity.Duration = lrw.Duration()

				if ratelimit.HandleRequest(activityHash, lrw.Status()) {
					activity.Status = http.StatusTooManyRequests
				}
			}

			go tracking.SaveActivity(activity, app)
		}
	})
}
