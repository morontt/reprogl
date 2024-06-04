package transliterator

import "testing"

func TestItShouldTransliterateGeneral(t *testing.T) {
	cases := map[string]string{
		"80 km/h":            "80 km/h",
		"дом":                "dom",
		"\u1eff":             "",
		"Александр Харченко": "Aleksandr Kharchenko",
		"Одесса Онищенко":    "Odessa Onishchenko",
		"Рыбатекст используется дизайнерами": "Rybatekst ispolzuetsia dizainerami",
		"Ёжик в тумане":         "Ezhik v tumane",
		"генератор бредотекста": "generator bredoteksta",
		"Зюзин Илья":            "Ziuzin Ilia",
		"Первый подъезд":        "Pervyi podieezd",
	}

	for text, expected := range cases {
		actual := Transliterate(text)
		if actual != expected {
			t.Errorf("Transliteration error: got %s; want %s", actual, expected)
		}
	}
}
