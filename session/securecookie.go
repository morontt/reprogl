package session

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"hash"
)

type SecureCookie struct {
	maxLength int
	encoded   string
	sz        serializer

	hashKey  []byte
	hashFunc func() hash.Hash
}

func NewSecureCookie(hashKey string) *SecureCookie {
	return &SecureCookie{
		maxLength: 4096,
		sz:        jsonEncoder{},
		hashKey:   []byte(hashKey),
	}
}

func (sc *SecureCookie) encode(data internalData) error {
	var err error
	var b []byte

	if b, err = sc.sz.serialize(data); err != nil {
		return err
	}

	mac := createMac(b, sc.hashKey)
	b = append(b, mac...)

	b = encode(b)
	if sc.maxLength != 0 && len(b) > sc.maxLength {
		return EncodedValueTooLong
	}

	sc.encoded = string(b)

	return nil
}

func (sc *SecureCookie) decode(value string) (internalData, error) {
	var err error
	var b []byte
	var data internalData

	b, err = decode(value)
	if err != nil {
		return data, err
	}

	if err = verifyMac(b, sc.hashKey); err != nil {
		return data, err
	}

	return sc.sz.deserialize(b)
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

func decode(value string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(value)
}

// createMac creates a message authentication code.
func createMac(value, key []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(value)

	return h.Sum(nil)
}

func verifyMac(value, key []byte) error {
	mac := createMac(value[:len(value)-32], key)
	if subtle.ConstantTimeCompare(value[len(value)-32:], mac) == 1 {
		return nil
	}

	return ErrMacInvalid
}
