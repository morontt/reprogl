package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"xelbot.com/reprogl/container"
)

func Recover(next http.Handler, app *container.Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil && rvr != http.ErrAbortHandler {
				err := errors.New(fmt.Sprintf("%v", rvr))
				app.ServerError(w, err)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
