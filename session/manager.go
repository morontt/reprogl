package session

import (
	"context"
	"errors"
	"net/http"

	"xelbot.com/reprogl/security"
)

const (
	CookieName = "session"
	CtxKey     = "session.ctx.key"
	CsrfCookie = "csrf_token"

	IdentityKey = "identity"

	VarnishSessionHeader = "X-Varnish-Session"
)

var (
	DecodeError         = errors.New("session: decode error")
	EncodedValueTooLong = errors.New("session: the encoded value is too long")
)

type CookieInterface interface {
	Name() string
	Path() string
	Value() string
	Persist() bool
}

func FromRequest(r *http.Request) (*Data, bool) {
	return newData(), true
}

func FromContext(ctx context.Context) *Data {
	c, ok := ctx.Value(CtxKey).(*Data)
	if !ok {
		panic("session: no data in context")
	}

	return c
}

func GetString(ctx context.Context, key string) (string, bool) {
	data := FromContext(ctx)

	data.mu.RLock()
	defer data.mu.RUnlock()

	if value, ok1 := data.values[key]; ok1 {
		if val, ok2 := value.(string); ok2 {
			return val, true
		}
	}

	return "", false
}

func Put(ctx context.Context, key string, value any) {
	data := FromContext(ctx)

	data.mu.Lock()
	data.values[key] = value
	data.status = Modified
	data.mu.Unlock()
}

func HasIdentity(ctx context.Context) bool {
	data := FromContext(ctx)

	data.mu.RLock()
	defer data.mu.RUnlock()

	if raw, exists := data.values[IdentityKey]; exists {
		_, ok := raw.(security.Identity)

		return ok
	}

	return false
}

func GetIdentity(ctx context.Context) (security.Identity, bool) {
	data := FromContext(ctx)

	data.mu.RLock()
	defer data.mu.RUnlock()

	if raw, exists := data.values[IdentityKey]; exists {
		if identity, ok := raw.(security.Identity); ok {
			return identity, true
		}
	}

	return security.Identity{}, false
}

func ClearIdentity(ctx context.Context) {
	data := FromContext(ctx)

	data.mu.Lock()
	defer data.mu.Unlock()

	_, exists := data.values[IdentityKey]
	if !exists {
		return
	}

	delete(data.values, IdentityKey)
	data.status = Modified
}
