package middlewares

import (
	"net/http"
	"time"
	"xelbot.com/reprogl/container"
)

func AccessLog(next http.Handler, app *container.Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		addr := r.Header.Get("X-Real-IP")
		if addr == "" {
			addr = r.Header.Get("X-Forwarded-For")
			if addr == "" {
				addr = r.RemoteAddr
			}
		}
		lrw := &logResponseWriter{w, 0}
		next.ServeHTTP(lrw, r)
		app.InfoLog.Printf("[%s] %s, %s %d %s\n", r.Method, addr, r.URL.Path, lrw.Status(), time.Since(start))
	})
}

type logResponseWriter struct {
	http.ResponseWriter
	status int
}

func (lrw *logResponseWriter) WriteHeader(statusCode int) {
	lrw.status = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

func (lrw *logResponseWriter) Status() int {
	if lrw.status == 0 {
		return http.StatusOK
	}

	return lrw.status
}
