package session

import (
	"net/http"
	"time"
)

type ResponseWriter struct {
	http.ResponseWriter
	sessionData *Data
	written     bool
}

func (sw *ResponseWriter) SetSessionData(d *Data) {
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
			expiry := time.Now().Add(14 * 24 * time.Hour)

			secureCookie = NewSecureCookie()
			err := secureCookie.Encode(sw.sessionData.values)
			if err != nil {
				panic(err)
			}

			writeCookie(sw, secureCookie, expiry)
		case Destroyed:
			secureCookie = NewSecureCookie()
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
