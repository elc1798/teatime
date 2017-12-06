package encode

import (
	json "encoding/json"
	"errors"
)

type ConnectionRequestSerializer struct{}

func (this *ConnectionRequestSerializer) Serialize(v interface{}) ([]byte, error) {
	obj, ok := v.(ConnectionRequestPayload)
	if !ok {
		return nil, errors.New("Invalid input")
	}
	return json.Marshal(obj)
}

func (this *ConnectionRequestSerializer) Deserialize(v []byte) (interface{}, error) {
	var data ConnectionRequestPayload
	if err := json.Unmarshal(v, &data); err != nil {
		return nil, err
	}

	return data, nil
}
