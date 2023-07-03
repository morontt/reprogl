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
