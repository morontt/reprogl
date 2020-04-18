package middlewares

import (
	"net/http"
	"time"
	"xelbot.com/reprogl/handlers"
)

func AccessLog(next http.Handler, app *handlers.Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		app.InfoLog.Printf("[%s] %s, %s %s\n", r.Method, r.RemoteAddr, r.URL.Path, time.Since(start))
	})
}
