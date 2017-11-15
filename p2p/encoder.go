package p2p

import (
	json "encoding/json"
)

type changedFileList struct {
	Filenames []string `json:"filenames"`
}

// Write interface for serializer
type Serializer interface {
	Serialize(v interface{}) ([]byte, error)
	Deserialize([]byte) (interface{}, error)
}

type ChangedFileListSerializer struct{}

func (this *ChangedFileListSerializer) Serialize(v interface{}) ([]byte, error) {
	list := v.([]string)
	data := changedFileList{Filenames: list}
	return json.Marshal(data)
}

func (this *ChangedFileListSerializer) Deserialize(v []byte) (interface{}, error) {
	var data changedFileList
	if err := json.Unmarshal(v, &data); err != nil {
		return nil, err
	}

	return data.Filenames, nil
}
