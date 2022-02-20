package middlewares

import (
	"net/http"
	pkghttp "xelbot.com/reprogl/http"
)

func ResponseWrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := &pkghttp.Response{ResponseWriter: w}

		next.ServeHTTP(lrw, r)
	})
}
