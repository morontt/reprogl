package models

import (
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
	FingerPrint  string
}

func (a *Activity) IsBot() bool {
	return isBot(a.UserAgent)
}
