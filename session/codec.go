package session

import (
	"bytes"
	"encoding/json"

	"xelbot.com/reprogl/security"
)

type serializer interface {
	serialize(src internalData) ([]byte, error)
	deserialize(src []byte) (internalData, error)
}

type jsonEncoder struct{}

func (e jsonEncoder) serialize(src internalData) ([]byte, error) {
	aux := struct {
		Identity security.Identity      `json:"a"`
		Values   map[string]interface{} `json:"v,omitempty"`
		Deadline deadline               `json:"d"`
	}{
		Identity: src.identity,
		Values:   src.values,
		Deadline: src.deadline,
	}

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(&aux); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (e jsonEncoder) deserialize(src []byte) (internalData, error) {
	aux := struct {
		Identity security.Identity      `json:"a"`
		Values   map[string]interface{} `json:"v,omitempty"`
		Deadline deadline               `json:"d"`
	}{}

	if err := json.NewDecoder(bytes.NewReader(src)).Decode(&aux); err != nil {
		return internalData{}, DecodeError
	}

	return internalData{
		identity: aux.Identity,
		values:   aux.Values,
		deadline: aux.Deadline,
	}, nil
}
