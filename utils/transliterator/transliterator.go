package transliterator

import (
	"strings"
	"unicode"
)

func Transliterate(text string) string {
	var replacement strings.Builder

	runes := []rune(text)
	for _, symbol := range runes {
		if symbol < unicode.MaxASCII {
			replacement.WriteString(string(symbol))

			continue
		}

		if symbol > 0x399 {
			idx := int(symbol - 0x400)
			if idx < len(x004) {
				replacement.WriteString(x004[idx])
			}
		}
	}

	return replacement.String()
}
