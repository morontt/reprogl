package session

import (
	"encoding/base64"
)

type SecureCookie struct {
	maxLength int
	encoded   string
	sz        Serializer
}

func NewSecureCookie() *SecureCookie {
	return &SecureCookie{
		maxLength: 4096,
		sz:        JSONEncoder{},
	}
}

func (sc *SecureCookie) Encode(value any) error {
	var err error
	var b []byte

	if b, err = sc.sz.Serialize(value); err != nil {
		return err
	}
	b = encode(b)

	if sc.maxLength != 0 && len(b) > sc.maxLength {
		return EncodedValueTooLong
	}

	sc.encoded = string(b)

	return nil
}

func (*SecureCookie) Name() string {
	return CookieName
}

func (*SecureCookie) Path() string {
	return "/"
}

func (*SecureCookie) Persist() bool {
	return true
}

func (sc *SecureCookie) Value() string {
	return sc.encoded
}

func encode(value []byte) []byte {
	encoded := make([]byte, base64.URLEncoding.EncodedLen(len(value)))
	base64.URLEncoding.Encode(encoded, value)

	return encoded
}
