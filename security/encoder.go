package security

import (
	"crypto/sha512"
	"fmt"
)

const iterations = 4000

func EncodePassword(password string, salt string) string {
	salted := password + "{" + salt + "}"
	h := sha512.New384()

	h.Write([]byte(salted))
	digest := h.Sum(nil)

	for i := 1; i < iterations; i++ {
		h.Reset()
		h.Write(digest)
		h.Write([]byte(salted))
		digest = h.Sum(nil)
	}

	return fmt.Sprintf("%x", digest)
}
