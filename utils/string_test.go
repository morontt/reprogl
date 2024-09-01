package utils

import (
	"strconv"
	"testing"
)

func TestEllipticalTruncate(t *testing.T) {
	tests := []struct {
		text   string
		want   string
		maxLen int
	}{
		{
			text:   "a",
			want:   "a",
			maxLen: 3,
		},
		{
			text:   "123",
			want:   "123",
			maxLen: 3,
		},
		{
			text:   "1 2 3 4 5 6",
			want:   "1 2 ...",
			maxLen: 7,
		},
		{
			text:   "Привет, безумный мир",
			want:   "Привет, безум...",
			maxLen: 16,
		},
	}

	for idx, item := range tests {
		t.Run(strconv.Itoa(idx), func(t *testing.T) {
			res := EllipticalTruncate(item.text, item.maxLen)

			if res != item.want {
				t.Errorf("%s : got %s; want %s", item.text, res, item.want)
			}
		})
	}
}
