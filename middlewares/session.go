package middlewares

import (
	"context"
	"log"
	"net/http"

	"xelbot.com/reprogl/session"
)

func Session(next http.Handler, logger *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionData := session.FromRequest(r, logger)

		ctx := r.Context()
		ctx = context.WithValue(ctx, session.CtxKey, sessionData)

		sw := &session.ResponseWriter{ResponseWriter: w}
		sw.SetSessionData(sessionData)

		next.ServeHTTP(sw, r.WithContext(ctx))

		sw.CheckAndWrite()
	})
}
