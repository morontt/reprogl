package hashid

import (
	"errors"
	"strconv"
	"testing"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		hash   string
		id     int
		isMale bool
		isUser bool
	}{
		{
			hash:   "ZQD5TM",
			id:     27,
			isUser: false,
			isMale: true,
		},
		{
			hash:   "04RETW",
			id:     27,
			isUser: false,
			isMale: false,
		},
		{
			hash:   "0WMMUN",
			id:     48,
			isUser: true,
			isMale: true,
		},
	}

	for idx, item := range tests {
		t.Run(strconv.Itoa(idx), func(t *testing.T) {
			model, err := Decode(item.hash)
			if err != nil {
				t.Errorf("decode error: %s", err.Error())
				return
			}

			if model.ID != item.id {
				t.Errorf("%s : ID got %d; want %d", item.hash, model.ID, item.id)
			}
			if model.IsUser() != item.isUser {
				t.Errorf("%s : wrong isUser detection", item.hash)
			}
			if model.IsMale() != item.isMale {
				t.Errorf("%s : wrong isMale detection", item.hash)
			}
		})
	}
}

func TestDecodeWithWrongOptions(t *testing.T) {
	tests := []string{"NVMW17", "XR5LU6", "ZDMLHM", "4R6HQ3"}

	for idx, item := range tests {
		t.Run(strconv.Itoa(idx), func(t *testing.T) {
			_, err := Decode(item)
			if err != nil && !errors.Is(err, WrongOptions) {
				t.Errorf("decode error: %s", err.Error())
			}
		})
	}
}
