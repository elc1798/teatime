package encode

import (
	"errors"
)

type ConnectionRequestSerializer struct{}

func (this *ConnectionRequestSerializer) Serialize(v interface{}) ([]byte, error) {
	obj, ok := v.(ConnectionRequestPayload)
	if !ok {
		return nil, errors.New("Invalid input")
	}
	return []byte(obj), nil // We know that it's a string internally
}

func (this *ConnectionRequestSerializer) Deserialize(v []byte) (interface{}, error) {
	return ConnectionRequestPayload(string(v)), nil
}
