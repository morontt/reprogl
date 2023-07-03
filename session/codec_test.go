package session

import (
	"fmt"
	"testing"
)

var testData = []map[string]any{
	{"key_1": "ABC", "key_2": 38},
	{"idx": []int{1, 2, 3}},
}

func TestJSONSerialization(t *testing.T) {
	var (
		codec      JSONEncoder
		serialized []byte
		err        error
	)

	for _, value := range testData {
		if serialized, err = codec.Serialize(value); err != nil {
			t.Error(err)
		} else {
			deserialized := make(map[string]any)
			if err = codec.Deserialize(serialized, &deserialized); err != nil {
				t.Error(err)
			}
			if fmt.Sprintf("%v", deserialized) != fmt.Sprintf("%v", value) {
				t.Errorf("Expected %v, got %v.", value, deserialized)
			}
		}
	}
}

func TestGobSerialization(t *testing.T) {
	var (
		codec      GobEncoder
		serialized []byte
		err        error
	)

	for _, value := range testData {
		if serialized, err = codec.Serialize(value); err != nil {
			t.Error(err)
		} else {
			deserialized := make(map[string]any)
			if err = codec.Deserialize(serialized, &deserialized); err != nil {
				t.Error(err)
			}
			if fmt.Sprintf("%v", deserialized) != fmt.Sprintf("%v", value) {
				t.Errorf("Expected %v, got %v.", value, deserialized)
			}
		}
	}
}
