package models

import (
	"crypto/md5"
	"fmt"
	"net"
	"time"
)

type Activity struct {
	Addr         net.IP
	UserAgent    string
	RequestedURI string
	Time         time.Time
	IsCDN        bool
	Status       int
}

func (a *Activity) AgentHash() string {
	hash := md5.New()
	hash.Write([]byte(a.UserAgent))

	return fmt.Sprintf("%x", hash.Sum(nil))
}
