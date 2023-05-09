package security

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"time"
	"xelbot.com/reprogl/container"
)

func GetWSSEHeader() (string, string) {
	cnf := container.GetConfig()
	created := time.Now().Format(time.RFC3339)

	nonce := make([]byte, 12)
	_, err := rand.Read(nonce)
	if err != nil {
		panic(err)
	}

	h := sha1.New()
	h.Write(nonce)
	h.Write([]byte(created))
	h.Write([]byte(cnf.BackendApiWsseKey))

	return "X-WSSE", fmt.Sprintf(
		`UsernameToken Username="%s", PasswordDigest="%s", Nonce="%s", Created="%s"`,
		cnf.BackendApiUser,
		base64.StdEncoding.EncodeToString(h.Sum(nil)),
		base64.StdEncoding.EncodeToString(nonce),
		created,
	)
}
