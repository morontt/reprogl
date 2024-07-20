package middlewares

import (
	"net/http"

	"xelbot.com/reprogl/container"
	pkghttp "xelbot.com/reprogl/http"
)

func AccessLog(next http.Handler, app *container.Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		addr := container.RealRemoteAddress(r)

		next.ServeHTTP(w, r)
		lrw, ok := w.(pkghttp.LogResponseWriter)
		if ok {
			app.InfoLog.Printf("[%s] %s, %s %d %s\n", r.Method, addr, r.URL.RequestURI(), lrw.Status(), lrw.Duration())
		} else {
			app.InfoLog.Printf("[%s] %s, %s\n", r.Method, addr, r.URL.RequestURI())
		}
	})
}
