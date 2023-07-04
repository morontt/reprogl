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

func FromRequest(r *http.Request) (*Store, bool) {
	return newStore(), true
}

func FromContext(ctx context.Context) *Store {
	c, ok := ctx.Value(CtxKey).(*Store)
	if !ok {
		panic("session: no data in context")
	}

	return c
}

func Put(ctx context.Context, key string, value any) {
	store := FromContext(ctx)

	store.mu.Lock()
	store.data.values[key] = value
	store.status = Modified
	store.mu.Unlock()
}

func Has(ctx context.Context, key string) bool {
	store := FromContext(ctx)

	store.mu.RLock()
	defer store.mu.RUnlock()

	if raw, exists := store.data.values[key]; exists {
		_, ok := raw.(security.Identity)

		return ok
	}

	return false
}

// TODO rework to generics
func GetString(ctx context.Context, key string) (string, bool) {
	store := FromContext(ctx)

	store.mu.RLock()
	defer store.mu.RUnlock()

	if raw, exists := store.data.values[key]; exists {
		if val, ok2 := raw.(string); ok2 {
			return val, true
		}
	}

	return "", false
}

func Remove(ctx context.Context, key string) {
	store := FromContext(ctx)

	store.mu.Lock()
	defer store.mu.Unlock()

	_, exists := store.data.values[key]
	if !exists {
		return
	}

	delete(store.data.values, key)
	store.status = Modified
}

func HasIdentity(ctx context.Context) (result bool) {
	store := FromContext(ctx)

	store.mu.RLock()
	result = !store.data.identity.IsZero()
	store.mu.RUnlock()

	return
}

func GetIdentity(ctx context.Context) (security.Identity, bool) {
	store := FromContext(ctx)

	store.mu.RLock()
	defer store.mu.RUnlock()

	return store.data.identity, !store.data.identity.IsZero()
}

func SetIdentity(ctx context.Context, identity security.Identity) {
	store := FromContext(ctx)

	store.mu.Lock()
	store.data.identity = identity
	store.status = Modified
	store.mu.Unlock()
}

func ClearIdentity(ctx context.Context) {
	store := FromContext(ctx)

	store.mu.Lock()
	defer store.mu.Unlock()

	if !store.data.identity.IsZero() {
		return
	}

	store.data.identity = security.Identity{}
	store.status = Modified
}
