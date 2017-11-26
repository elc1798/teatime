package encode

import (
	json "encoding/json"
	"errors"
)

type PingSerializer struct{}

func (this *PingSerializer) Serialize(v interface{}) ([]byte, error) {
	obj, ok := v.(PingPayload)
	if !ok {
		return nil, errors.New("Invalid input")
	}
	return json.Marshal(obj)
}

func (this *PingSerializer) Deserialize(v []byte) (interface{}, error) {
	var data PingPayload
	if err := json.Unmarshal(v, &data); err != nil {
		return nil, err
	}

	return data, nil
}
