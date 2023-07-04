package session

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"xelbot.com/reprogl/security"
)

type serializer interface {
	serialize(src internalData) ([]byte, error)
	deserialize(src []byte) (internalData, error)
}

type gobEncoder struct{}
type jsonEncoder struct{}

func (e gobEncoder) serialize(src internalData) ([]byte, error) {
	aux := struct {
		Identity security.Identity
		Values   map[string]interface{}
	}{
		Identity: src.identity,
		Values:   src.values,
	}

	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(&aux); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (e gobEncoder) deserialize(src []byte) (internalData, error) {
	aux := struct {
		Identity security.Identity
		Values   map[string]interface{}
	}{}

	if err := gob.NewDecoder(bytes.NewBuffer(src)).Decode(&aux); err != nil {
		return internalData{}, DecodeError
	}

	return internalData{
		identity: aux.Identity,
		values:   aux.Values,
	}, nil
}

func (e jsonEncoder) serialize(src internalData) ([]byte, error) {
	aux := struct {
		Identity security.Identity
		Values   map[string]interface{}
	}{
		Identity: src.identity,
		Values:   src.values,
	}

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(&aux); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (e jsonEncoder) deserialize(src []byte) (internalData, error) {
	aux := struct {
		Identity security.Identity
		Values   map[string]interface{}
	}{}

	if err := json.NewDecoder(bytes.NewReader(src)).Decode(&aux); err != nil {
		return internalData{}, DecodeError
	}

	return internalData{
		identity: aux.Identity,
		values:   aux.Values,
	}, nil
}
