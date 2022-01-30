package middlewares

import (
	"fmt"
	"net/http"
	"runtime"
	"xelbot.com/reprogl/container"
	pkghttp "xelbot.com/reprogl/http"
)

func ResponseWrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := &pkghttp.Response{ResponseWriter: w}
		addXPoweredBy(lrw)

		next.ServeHTTP(lrw, r)
	})
}

func addXPoweredBy(w http.ResponseWriter) {
	w.Header().Set("X-Powered-By", fmt.Sprintf(
		"Reprogl/%s (%s)",
		container.GitRevision,
		runtime.Version()))
}
