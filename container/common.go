package container

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"strings"
)

func RealRemoteAddress(r *http.Request) string {
	addr := r.Header.Get("X-Real-IP")
	if addr == "" {
		addr = r.Header.Get("X-Forwarded-For")
		if addr == "" {
			addr = r.RemoteAddr
		}
	}

	return addr
}

func IsCDN(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Via"), "BunnyCDN")
}

func MD5(s string) string {
	hash := md5.Sum([]byte(s))

	return fmt.Sprintf("%x", hash)
}
