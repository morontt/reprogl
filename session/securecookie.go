package session

import (
	"crypto/cipher"
	"crypto/des"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
)

type SecureCookie struct {
	maxLength int
	encoded   string
	sz        serializer

	// For testing purposes
	skipExpiration bool

	hashKey []byte
	block   cipher.Block
}

func NewSecureCookie(hashKey, cipherKey string) *SecureCookie {
	h := sha256.New()
	h.Write([]byte(cipherKey))
	cipherKeyHash := h.Sum(nil)
	block, _ := des.NewTripleDESCipher(cipherKeyHash[:24])

	return &SecureCookie{
		maxLength: 4096,
		sz:        jsonEncoder{},
		hashKey:   []byte(hashKey),
		block:     block,
	}
}

func (sc *SecureCookie) encode(data internalData) error {
	var err error
	var b []byte

	if b, err = sc.sz.serialize(data); err != nil {
		return err
	}

	if sc.block != nil {
		if b, err = encrypt(sc.block, b); err != nil {
			return err
		}
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

	if sc.block != nil {
		if b, err = decrypt(sc.block, b); err != nil {
			return data, err
		}
	}

	data, err = sc.sz.deserialize(b)
	if err != nil {
		return data, err
	}

	if !sc.skipExpiration && data.deadline.IsExpired() {
		return data, Expired
	}

	return data, nil
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

// For testing purposes
func (sc *SecureCookie) ignoreExpiration() {
	sc.skipExpiration = true
}

func encrypt(block cipher.Block, value []byte) ([]byte, error) {
	iv := make([]byte, block.BlockSize())
	_, err := rand.Read(iv)
	if err != nil {
		return nil, EncryptionError
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(value, value)

	return append(iv, value...), nil
}

func decrypt(block cipher.Block, value []byte) ([]byte, error) {
	size := block.BlockSize()
	if len(value) > size {
		iv := value[:size]
		value = value[size:]
		stream := cipher.NewCTR(block, iv)
		stream.XORKeyStream(value, value)
		return value, nil
	}

	return nil, DecryptionError
}
