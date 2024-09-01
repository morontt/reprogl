package utils

func EllipticalTruncate(text string, maxLen int) string {
	runes := []rune(text)
	if len(runes) <= maxLen {
		return text
	}
	if maxLen < 3 {
		maxLen = 3
	}

	return string(runes[0:maxLen-3]) + "..."
}
