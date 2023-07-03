package middlewares

import (
	"context"
	"net/http"

	"xelbot.com/reprogl/session"
)

func Session(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionData, _ := session.FromRequest(r)

		ctx := r.Context()
		ctx = context.WithValue(ctx, session.CtxKey, sessionData)

		sw := &session.ResponseWriter{ResponseWriter: w}
		sw.SetSessionData(sessionData)

		next.ServeHTTP(sw, r.WithContext(ctx))

		sw.CheckAndWrite()
	})
}
