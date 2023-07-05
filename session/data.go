package session

import (
	"encoding/json"
	"sync"
	"time"

	"xelbot.com/reprogl/security"
)

type Status uint8

const (
	Unmodified Status = iota
	Modified
	Destroyed
)

type deadline time.Time

type Store struct {
	status Status
	data   internalData
	mu     sync.RWMutex
}

type internalData struct {
	values   map[string]interface{}
	identity security.Identity
	deadline deadline
}

func newStore() *Store {
	return &Store{
		status: Unmodified,
		data: internalData{
			values: make(map[string]interface{}),
		},
	}
}

func newStoreWithData(d internalData) *Store {
	return &Store{
		status: Unmodified,
		data:   d,
	}
}

func (d *deadline) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}
	*d = deadline(t)

	return nil
}

func (d deadline) MarshalJSON() ([]byte, error) {
	t := time.Time(d)

	return json.Marshal(t.Format(time.RFC3339))
}

func (s *Store) setModified() {
	if s.status != Destroyed {
		s.status = Modified
	}
}

func (d deadline) IsExpired() bool {
	t := time.Time(d)

	return t.Before(time.Now())
}
