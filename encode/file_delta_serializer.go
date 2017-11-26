package encode

import (
	json "encoding/json"
	"errors"
)

type FileDeltasSerializer struct{}

func (this *FileDeltasSerializer) Serialize(v interface{}) ([]byte, error) {
	obj, ok := v.(FileDeltasPayload)
	if !ok {
		return nil, errors.New("Invalid input")
	}
	return json.Marshal(obj)
}

func (this *FileDeltasSerializer) Deserialize(v []byte) (interface{}, error) {
	var data FileDeltasPayload
	if err := json.Unmarshal(v, &data); err != nil {
		return nil, err
	}

	return data, nil
}
