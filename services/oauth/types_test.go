package oauth

import (
	"encoding/json"
	"strconv"
	"testing"
)

func TestGender_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		jsonString string
		want       Gender
	}{
		{
			jsonString: `{"a": "b", "sex": "male"}`,
			want:       Male,
		},
		{
			jsonString: `{"a": "b", "sex": "female"}`,
			want:       Female,
		},
		{
			jsonString: `{"a": "b", "sex": null}`,
			want:       Unknown,
		},
		{
			jsonString: `{"a": "b"}`,
			want:       Unknown,
		},
	}

	for idx, item := range tests {
		t.Run(strconv.Itoa(idx), func(t *testing.T) {
			var res = struct {
				Gender Gender `json:"sex"`
			}{
				Gender: Unknown,
			}

			err := json.Unmarshal([]byte(item.jsonString), &res)
			if err != nil {
				t.Error(err)
			}

			if res.Gender != item.want {
				t.Errorf("%s : got %s; want %s", item.jsonString, res.Gender, item.want)
			}
		})
	}
}
