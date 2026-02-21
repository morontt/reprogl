package models

import (
	"net"
	"strings"
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
	LocationID   int
	Duration     time.Duration
	Method       string
}

func (a *Activity) IsBot() bool {
	return isBot(a.UserAgent)
}

func (a *Activity) IsInternalRequest() bool {
	return strings.HasPrefix(a.RequestedURI, "/_fragment/")
}
