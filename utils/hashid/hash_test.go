package hashid

import (
	"strconv"
	"testing"
)

func TestInternalDecode(t *testing.T) {
	tests := []struct {
		hash    string
		id      int
		options Option
	}{
		{
			hash:    "ZQD5TM",
			id:      27,
			options: Commentator | Male,
		},
		{
			hash:    "04RETW",
			id:      27,
			options: Commentator | Female,
		},
		{
			hash:    "0WMMUN",
			id:      48,
			options: User | Male,
		},
	}

	for idx, item := range tests {
		t.Run(strconv.Itoa(idx), func(t *testing.T) {
			ids, err := decode(item.hash)
			if err != nil {
				t.Errorf("decode error: %s", err.Error())
				return
			}

			if len(ids) != 2 {
				t.Errorf("%s : incorrect ids len: %d", item.hash, len(ids))
				return
			}

			if ids[0] != item.id {
				t.Errorf("%s : ID got %d; want %d", item.hash, ids[0], item.id)
			}

			if ids[1] != int(item.options) {
				t.Errorf("%s : Options got %d; want %d", item.hash, ids[1], item.options)
			}
		})
	}
}
