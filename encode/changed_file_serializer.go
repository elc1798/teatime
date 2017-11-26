package encode

import (
	json "encoding/json"
	"errors"
)

type ChangedFileListSerializer struct{}

func (this *ChangedFileListSerializer) Serialize(v interface{}) ([]byte, error) {
	obj, ok := v.(ChangedFileListPayload)
	if !ok {
		return nil, errors.New("Invalid input")
	}
	return json.Marshal(obj)
}

func (this *ChangedFileListSerializer) Deserialize(v []byte) (interface{}, error) {
	var data ChangedFileListPayload
	if err := json.Unmarshal(v, &data); err != nil {
		return nil, err
	}

	return data, nil
}
