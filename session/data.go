package session

import "sync"

type Status uint8

const (
	Unmodified Status = iota
	Modified
	Destroyed
)

type Data struct {
	status Status
	values map[string]interface{}
	mu     sync.RWMutex
}

func newData() *Data {
	return &Data{
		status: Unmodified,
		values: make(map[string]interface{}),
	}
}
