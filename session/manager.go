package session

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/security"
)

const (
	CookieName = "session"
	CtxKey     = "session.ctx.key"
	CsrfCookie = "csrf_token"

	FlashErrorKey   = "fe"
	FlashSuccessKey = "fs"

	varnishSessionHeader = "X-Varnish-Session"

	maxAge = 14 * 24 * time.Hour
)

var (
	DecodeError         = errors.New("session: decode error")
	EncodedValueTooLong = errors.New("session: the encoded value is too long")
	ErrMacInvalid       = errors.New("session: the HMAC is not valid")
	Expired             = errors.New("session: expired")
)

type CookieInterface interface {
	Name() string
	Path() string
	Value() string
	Persist() bool
}

func FromRequest(r *http.Request, logger *log.Logger) *Store {
	var withError bool
	requestData := r.Header.Get(varnishSessionHeader)
	if len(requestData) > 0 {
		secureCookie := NewSecureCookie(container.GetConfig().SessionHashKey)
		data, err := secureCookie.decode(requestData)
		if err == nil {
			return newStoreWithData(data)
		} else {
			withError = true
			logger.Printf("[AUTH] session: %s error: %s\n", requestData, err.Error())
		}
	}

	store := newStore()
	if withError {
		store.mu.Lock()
		store.status = Destroyed
		store.mu.Unlock()
	}

	return store
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
	store.setModified()
	store.mu.Unlock()
}

func Has(ctx context.Context, key string) bool {
	store := FromContext(ctx)

	store.mu.RLock()
	_, exists := store.data.values[key]
	store.mu.RUnlock()

	return exists
}

func Get[T any](ctx context.Context, key string) (T, bool) {
	store := FromContext(ctx)

	store.mu.RLock()
	defer store.mu.RUnlock()

	var result T
	var ok bool
	if raw, exists := store.data.values[key]; exists {
		if result, ok = raw.(T); ok {
			return result, true
		}
	}

	return result, false
}

func Pop[T any](ctx context.Context, key string) (T, bool) {
	store := FromContext(ctx)

	store.mu.Lock()
	defer store.mu.Unlock()

	var result T
	var ok bool
	if raw, exists := store.data.values[key]; exists {
		if result, ok = raw.(T); ok {
			delete(store.data.values, key)
			store.setModified()

			return result, true
		}
	}

	return result, false
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
	store.setModified()
}

func Destroy(ctx context.Context) {
	store := FromContext(ctx)

	store.mu.Lock()
	store.status = Destroyed
	store.mu.Unlock()
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
	store.setModified()
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
	store.setModified()
}
