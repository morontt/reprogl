package middlewares

import (
	"net/http"
	"strings"

	"xelbot.com/reprogl/container"
)

func CDN(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if container.IsCDN(r) {
			if r.URL.Path == "" || r.URL.Path == "/" {
				w.Header().Set("Content-Type", "text/plain")
				w.Write([]byte("This is static assets storage"))

				return
			} else if r.URL.Path != "/robots.txt" &&
				r.URL.Path != "/sitemap.xml" &&
				r.URL.Path != "/headers" &&
				!strings.HasPrefix(r.URL.Path, "/images/avatar/") {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)

				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
