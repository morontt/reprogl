package session

import (
	"net/http"
	"time"

	"xelbot.com/reprogl/container"
)

type ResponseWriter struct {
	http.ResponseWriter
	sessionData *Store
	written     bool
}

func (sw *ResponseWriter) SetSessionData(d *Store) {
	sw.sessionData = d
}

func (sw *ResponseWriter) Write(b []byte) (int, error) {
	sw.CheckAndWrite()

	return sw.ResponseWriter.Write(b)
}

func (sw *ResponseWriter) WriteHeader(code int) {
	sw.CheckAndWrite()

	sw.ResponseWriter.WriteHeader(code)
}

func (sw *ResponseWriter) CheckAndWrite() {
	if !sw.written {
		var secureCookie *SecureCookie
		switch sw.sessionData.status {
		case Modified:
			expiry := time.Now().Add(maxAge)

			secureCookie = NewSecureCookie(container.GetConfig().SessionHashKey)

			sw.sessionData.mu.Lock()
			sw.sessionData.data.deadline = deadline(expiry)
			err := secureCookie.encode(sw.sessionData.data)
			sw.sessionData.mu.Unlock()

			if err != nil {
				panic(err)
			}

			writeCookie(sw, secureCookie, expiry)
		case Destroyed:
			secureCookie = NewSecureCookie(container.GetConfig().SessionHashKey)
			writeCookie(sw, secureCookie, time.Time{})
		}

		if len(sw.Header().Values("Set-Cookie")) > 0 {
			sw.Header().Set("Cache-Control", `no-cache="Set-Cookie"`)
		}
	}

	sw.written = true
}

func writeCookie(w http.ResponseWriter, c CookieInterface, expiry time.Time) {
	cookie := &http.Cookie{
		Name:     c.Name(),
		Value:    c.Value(),
		Path:     c.Path(),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
	}

	if expiry.IsZero() {
		cookie.Expires = time.Unix(1, 0)
		cookie.MaxAge = -1
	} else if c.Persist() {
		cookie.Expires = time.Unix(expiry.Unix()+1, 0)
		cookie.MaxAge = int(time.Until(expiry).Seconds() + 1)
	}

	http.SetCookie(w, cookie)
}
