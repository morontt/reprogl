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
		expiry := time.Now().Add(14 * 24 * time.Hour)
		switch sw.sessionData.status {
		case Modified:
			writeSessionCookie(sw, sw.sessionData, expiry)
		case Destroyed:
			writeSessionCookie(sw, nil, time.Time{})
		}
	}

	sw.written = true
}

func writeSessionCookie(w http.ResponseWriter, d *Data, expiry time.Time) {
	secureCookie := NewSecureCookie()
	encoded, err := secureCookie.Encode(d.values)
	if err != nil {
		panic(err)
	}

	cookie := &http.Cookie{
		Name:     "session",
		Value:    encoded,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
	}

	if expiry.IsZero() {
		cookie.Expires = time.Unix(1, 0)
		cookie.MaxAge = -1
	} else {
		cookie.Expires = time.Unix(expiry.Unix()+1, 0)
		cookie.MaxAge = int(time.Until(expiry).Seconds() + 1)
	}

	w.Header().Set("Cache-Control", `no-cache="Set-Cookie"`)
	w.Header().Add("Set-Cookie", cookie.String())
}
