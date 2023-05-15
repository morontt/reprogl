package middlewares

import (
	"net/http"
	"time"
	"xelbot.com/reprogl/container"
	pkghttp "xelbot.com/reprogl/http"
)

func AccessLog(next http.Handler, app *container.Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		addr := container.RealRemoteAddress(r)

		next.ServeHTTP(w, r)
		lrw, ok := w.(pkghttp.LogResponseWriter)
		if ok {
			app.InfoLog.Printf("[%s] %s, %s %d %s\n", r.Method, addr, r.URL.Path, lrw.Status(), time.Since(start))
		} else {
			app.InfoLog.Printf("[%s] %s, %s %s\n", r.Method, addr, r.URL.Path, time.Since(start))
		}
	})
}
