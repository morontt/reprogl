package middlewares

import (
	"math/rand"
	"net/http"
)

type ClacksSet []string

func (c ClacksSet) Name() string {
	return "GNU " + c[rand.Intn(len(c))]
}

func (c ClacksSet) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Clacks-Overhead", c.Name())

		next.ServeHTTP(w, r)
	})
}

func Clacks() ClacksSet {
	data := []string{
		"Terry Pratchett",
		"Clive Sinclair",
	}

	return ClacksSet(data)
}
