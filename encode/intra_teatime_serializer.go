package encode

import (
	json "encoding/json"
	"errors"
)

const COMMAND_INIT_REPO = "init_repo"
const COMMAND_ADD_FILE = "add_file"
const COMMAND_LINK_PEER = "link_peer"

type intraTeatimePackage struct {
	Command   string   `json:"command"`
	Arguments []string `json:"arguments"`
}

type IntraTeatimeSerializer struct{}

func (this *IntraTeatimeSerializer) Serialize(v interface{}) ([]byte, error) {
	obj, ok := v.([]string)
	if !ok {
		return nil, errors.New("Invalid input")
	}

	// Serialize data: obj[0] is command, obj[1:] is arguments
	return json.Marshal(intraTeatimePackage{
		Command:   obj[0],
		Arguments: obj[1:],
	})
}

func (this *IntraTeatimeSerializer) Deserialize(v []byte) (interface{}, error) {
	var pack intraTeatimePackage
	if err := json.Unmarshal(v, &pack); err != nil {
		return nil, err
	}

	return append([]string{pack.Command}, pack.Arguments...), nil
}
