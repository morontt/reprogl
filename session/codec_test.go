package session

import (
	"fmt"
	"testing"

	"xelbot.com/reprogl/security"
)

var testData = []internalData{
	{
		identity: security.Identity{ID: 13, Username: "pupa"},
	},
	{
		identity: security.Identity{},
		values: map[string]any{
			"token": "random data",
			"width": 15,
		},
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
		if serialized, err = codec.serialize(value); err != nil {
			t.Error(err)
		} else {
			if deserialized, err = codec.deserialize(serialized); err != nil {
				t.Error(err)
			}
			if fmt.Sprintf("%+v", deserialized) != fmt.Sprintf("%+v", value) {
				t.Errorf("Expected %+v, got %+v.", value, deserialized)
			}
		}
	}
}

func TestGobSerialization(t *testing.T) {
	var (
		codec        gobEncoder
		serialized   []byte
		deserialized internalData
		err          error
	)

	for _, value := range testData {
		if serialized, err = codec.serialize(value); err != nil {
			t.Error(err)
		} else {
			if deserialized, err = codec.deserialize(serialized); err != nil {
				t.Error(err)
			}
			if fmt.Sprintf("%+v", deserialized) != fmt.Sprintf("%+v", value) {
				t.Errorf("Expected %+v, got %+v.", value, deserialized)
			}
		}
	}
}
