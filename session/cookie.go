package session

import (
	"net/http"
	"time"
)

type internalCookie struct {
	name, value, path string
	persist           bool
}

func WriteSessionCookie(w http.ResponseWriter, name, value, path string) {
	cookie := &internalCookie{
		name:  name,
		value: value,
		path:  path,
	}

	writeCookie(w, cookie, time.Now())
}

func WritePermanentCookie(w http.ResponseWriter, name, value, path string, expiry time.Time) {
	cookie := &internalCookie{
		name:    name,
		value:   value,
		path:    path,
		persist: true,
	}

	writeCookie(w, cookie, expiry)
}

// DeleteCookie note: Finally, to remove a cookie, the server returns a Set-Cookie header
// with an expiration date in the past. The server will be successful in removing the
// cookie only if the Path and the Domain attribute in the Set-Cookie header match the values
// used when the cookie was created.
//
// From: https://www.rfc-editor.org/rfc/rfc6265.html
func DeleteCookie(w http.ResponseWriter, name, path string) {
	cookie := &internalCookie{
		name:  name,
		value: "deleted",
		path:  path,
	}

	writeCookie(w, cookie, time.Time{})
}

func (c *internalCookie) Name() string {
	return c.name
}

func (c *internalCookie) Path() string {
	return c.path
}

func (c *internalCookie) Value() string {
	return c.value
}

func (c *internalCookie) Persist() bool {
	return c.persist
}
