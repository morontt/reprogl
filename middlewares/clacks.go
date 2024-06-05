package middlewares

import (
	"math/rand"
	"net/http"

	"xelbot.com/reprogl/utils/transliterator"
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
		"Николай Заманов",
		"Владлен Татарский",
		"Арсен Павлов",
		"Михаил Толстых",
		"Олесь Бузина",
		"Алексей Мозговой",
		"Robert Sheckley",
		"Robert Anson Heinlein",
	}

	tData := make(ClacksSet, len(data))
	for i, d := range data {
		tData[i] = transliterator.Transliterate(d)
	}

	return tData
}
