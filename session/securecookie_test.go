package session

import (
	"fmt"
	"testing"
)

func TestSecureCookie(t *testing.T) {
	var (
		deserialized internalData
		err          error
	)

	secureCookie := NewSecureCookie()
	for _, value := range testData {
		if err = secureCookie.encode(value); err != nil {
			t.Error(err)
		}

		serialized := secureCookie.Value()
		if deserialized, err = secureCookie.decode(serialized); err != nil {
			t.Error(err)
		}

		if fmt.Sprintf("%+v", deserialized) != fmt.Sprintf("%+v", value) {
			t.Errorf("Expected %+v, got %+v.", value, deserialized)
		}

	}
}
