package session

import (
	"context"
	"net/http"
)

const (
	CtxKey = "session.ctx.key"
)

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
	if value, ok2 := data.values[key]; ok2 {
		if val, ok3 := value.(string); ok3 {
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
