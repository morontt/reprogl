package session

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
)

type Serializer interface {
	Serialize(src interface{}) ([]byte, error)
	Deserialize(src []byte, dst interface{}) error
}

type GobEncoder struct{}
type JSONEncoder struct{}

func (e GobEncoder) Serialize(src any) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(src); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (e GobEncoder) Deserialize(src []byte, dst any) error {
	dec := gob.NewDecoder(bytes.NewBuffer(src))
	if err := dec.Decode(dst); err != nil {
		return DecodeError
	}

	return nil
}

func (e JSONEncoder) Serialize(src any) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	if err := enc.Encode(src); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (e JSONEncoder) Deserialize(src []byte, dst any) error {
	dec := json.NewDecoder(bytes.NewReader(src))
	if err := dec.Decode(dst); err != nil {
		return DecodeError
	}

	return nil
}
