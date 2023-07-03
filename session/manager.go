package session

import (
	"context"
	"errors"
	"net/http"
)

const (
	CookieName = "session"
	CtxKey     = "session.ctx.key"
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
