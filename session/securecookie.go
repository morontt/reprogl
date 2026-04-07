package session

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"hash"
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
	h := (hashFunc())()
	if h.Size() < 24 {
		panic(errors.New("session: invalid hash size"))
	}

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

	if b, err = verifyMac(b, sc.hashKey); err != nil {
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
	h := hmac.New(hashFunc(), key)
	h.Write(value)

	return h.Sum(nil)
}

func verifyMac(value, key []byte) ([]byte, error) {
	h := (hashFunc())()
	if len(value) <= h.Size() {
		return nil, ErrMacInvalid
	}

	data := value[:len(value)-h.Size()]
	mac := createMac(data, key)
	if subtle.ConstantTimeCompare(value[len(value)-h.Size():], mac) == 1 {
		return data, nil
	}

	return nil, ErrMacInvalid
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

	mode := cipher.NewCBCEncrypter(block, iv)

	plaintext := pad(value, block.BlockSize())
	ciphertext := make([]byte, len(plaintext))
	mode.CryptBlocks(ciphertext, plaintext)

	return append(iv, ciphertext...), nil
}

func decrypt(block cipher.Block, value []byte) ([]byte, error) {
	var err error
	size := block.BlockSize()
	if len(value) > size {
		iv := value[:size]
		value = value[size:]

		// input not full blocks
		if len(value)%size != 0 {
			return nil, DecryptionError
		}

		mode := cipher.NewCBCDecrypter(block, iv)
		plaintext := make([]byte, len(value))
		mode.CryptBlocks(plaintext, value)
		plaintext, err = unpad(plaintext, size)
		if err != nil {
			return nil, err
		}

		return plaintext, nil
	}

	return nil, DecryptionError
}

func hashFunc() func() hash.Hash {
	return sha256.New224
}

func pad(src []byte, blockSize int) []byte {
	padLen := blockSize - len(src)%blockSize
	padding := bytes.Repeat([]byte{byte(padLen)}, padLen)

	return append(src, padding...)
}

func unpad(src []byte, blockSize int) ([]byte, error) {
	padLen := int(src[len(src)-1])

	// slice bounds out of range
	if padLen < 1 || padLen > blockSize {
		return nil, DecryptionError
	}

	return src[:len(src)-padLen], nil
}
