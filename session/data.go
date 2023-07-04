package session

import (
	"sync"

	"xelbot.com/reprogl/security"
)

type Status uint8

const (
	Unmodified Status = iota
	Modified
	Destroyed
)

type Store struct {
	status Status
	data   internalData
	mu     sync.RWMutex
}

type internalData struct {
	values   map[string]interface{}
	identity security.Identity
}

func newStore() *Store {
	return &Store{
		status: Unmodified,
		data: internalData{
			values: make(map[string]interface{}),
		},
	}
}
