package session

import (
	"encoding/base64"
	"errors"
)

var (
	EncodedValueTooLong = errors.New("session: the encoded value is too long")
)

type SecureCookie struct {
	maxLength int
	sz        Serializer
}

func NewSecureCookie() *SecureCookie {
	return &SecureCookie{
		maxLength: 4096,
		sz:        JSONEncoder{},
	}
}

func (sc *SecureCookie) Encode(value any) (string, error) {
	var err error
	var b []byte

	if b, err = sc.sz.Serialize(value); err != nil {
		return "", err
	}
	b = encode(b)

	if sc.maxLength != 0 && len(b) > sc.maxLength {
		return "", EncodedValueTooLong
	}

	return string(b), nil
}

func encode(value []byte) []byte {
	encoded := make([]byte, base64.URLEncoding.EncodedLen(len(value)))
	base64.URLEncoding.Encode(encoded, value)

	return encoded
}
