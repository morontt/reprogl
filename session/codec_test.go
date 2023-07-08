package session

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"xelbot.com/reprogl/security"
)

type testDataItem struct {
	data       internalData
	jsonString string
}

var time1, _ = time.Parse(time.RFC822, "05 Jul 23 16:24 MSK")
var time2, _ = time.Parse(time.RFC822, "01 Feb 19 01:37 EET")

var testData = []testDataItem{
	{
		data: internalData{
			identity: security.Identity{ID: 13, Username: "pupa", Role: security.Admin},
			deadline: deadline(time1),
		},
		jsonString: `{"a":{"i":13,"u":"pupa","r":"admin"},"d":"2023-07-05T16:24:00+03:00"}`,
	},
	{
		data: internalData{
			identity: security.Identity{},
			values: map[string]any{
				"token": "random data",
				"width": 15,
			},
			deadline: deadline(time2),
		},
		jsonString: `{"a":{},"v":{"token":"random data","width":15},"d":"2019-02-01T02:37:00+03:00"}`,
	},
	{
		data: internalData{
			identity: security.Identity{ID: 7, Username: "lupa", Role: security.User},
			values: map[string]any{
				"abc": "zxc",
			},
		},
		jsonString: `{"a":{"i":7,"u":"lupa","r":"user"},"v":{"abc":"zxc"},"d":"0001-01-01T00:00:00Z"}`,
	},
}

func TestJSONSerialization(t *testing.T) {
	var (
		codec        jsonEncoder
		serialized   []byte
		deserialized internalData
		err          error
	)

	for _, value := range testData {
		if serialized, err = codec.serialize(value.data); err != nil {
			t.Error(err)
		} else {
			if deserialized, err = codec.deserialize(serialized); err != nil {
				t.Error(err)
			}
			if fmt.Sprintf("%+v", deserialized) != fmt.Sprintf("%+v", value.data) {
				t.Errorf("Expected %+v, got %+v", value.data, deserialized)
			}
		}
	}
}

func TestJSONString(t *testing.T) {
	var (
		codec      jsonEncoder
		serialized []byte
		str        string
		err        error
	)

	for _, value := range testData {
		if serialized, err = codec.serialize(value.data); err != nil {
			t.Error(err)
		} else {
			str = strings.TrimRight(string(serialized), "\n")
			if str != value.jsonString {
				t.Errorf("Expected %s, got %+v", value.jsonString, str)
			}
		}
	}
}
