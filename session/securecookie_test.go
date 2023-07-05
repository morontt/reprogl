package session

import (
	"errors"
	"fmt"
	"testing"
)

func TestSecureCookie(t *testing.T) {
	var (
		deserialized internalData
		err          error
	)

	secureCookie := NewSecureCookie("hash key")
	secureCookie.ignoreExpiration()
	for _, value := range testData {
		if err = secureCookie.encode(value.data); err != nil {
			t.Error(err)
		}

		serialized := secureCookie.Value()
		if deserialized, err = secureCookie.decode(serialized); err != nil {
			t.Error(err)
		}

		if fmt.Sprintf("%+v", deserialized) != fmt.Sprintf("%+v", value.data) {
			t.Errorf("Expected %+v, got %+v.", value.data, deserialized)
		}
	}
}

func TestInvalidHMAC(t *testing.T) {
	var (
		err error
	)

	secureCookie := NewSecureCookie("Lorem ipsum...")
	for _, value := range testData {
		if err = secureCookie.encode(value.data); err != nil {
			t.Error(err)
		}

		raw := []byte(secureCookie.Value())
		raw[1] = 'N'
		raw[17] = 'u'
		serialized := string(raw)
		if _, err = secureCookie.decode(serialized); err != nil {
			if !errors.Is(err, ErrMacInvalid) {
				t.Error(err)
			}
		} else {
			t.Errorf("Expected invalid HMAC for %+v", value.data)
		}
	}
}
