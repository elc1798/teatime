package encode

import (
	json "encoding/json"
	"errors"
)

const ACTION_CONNECT = "tt_conn_request"
const ACTION_PING = "ping"
const ACTION_FILE_LIST = "changed_f_list"
const ACTION_DELTAS = "deltas"

type routedPackage struct {
	Destination string `json:"destination"`
	Action      string `json:"action"`
	Data        string `json:"data"`
}

type InterTeatimeSerializer struct{}

func serializerFromAction(action string) (Serializer, error) {
	switch action {
	case ACTION_CONNECT:
		return &ConnectionRequestSerializer{}, nil
	case ACTION_PING:
		return &PingSerializer{}, nil
	case ACTION_FILE_LIST:
		return &ChangedFileListSerializer{}, nil
	case ACTION_DELTAS:
		return &FileDeltasSerializer{}, nil
	default:
		return nil, errors.New("Invalid action received!")
	}
}

func (this *InterTeatimeSerializer) Serialize(v interface{}) ([]byte, error) {
	obj, ok := v.(TeatimeMessage)
	if !ok {
		return nil, errors.New("Invalid input")
	}

	// Serialize "Data" based on action type
	payloadSerializer, err := serializerFromAction(obj.Action)
	if err != nil {
		return nil, err
	}

	data, err := payloadSerializer.Serialize(obj.Payload)
	if err != nil {
		return nil, err
	}

	return json.Marshal(routedPackage{
		Destination: obj.Recipient,
		Action:      obj.Action,
		Data:        string(data),
	})
}

func (this *InterTeatimeSerializer) Deserialize(v []byte) (interface{}, error) {
	var pack routedPackage
	if err := json.Unmarshal(v, &pack); err != nil {
		return nil, err
	}

	newMessage := TeatimeMessage{
		Recipient: pack.Destination,
		Action:    pack.Action,
		Payload:   nil,
	}

	// Deserialize "Data" based on action type
	payloadSerializer, err := serializerFromAction(pack.Action)
	if err != nil {
		return nil, err
	}

	payload, err := payloadSerializer.Deserialize([]byte(pack.Data))
	if err != nil {
		return nil, err
	}

	newMessage.Payload = payload

	return newMessage, nil
}

type TeatimeMessage struct {
	Recipient string
	Action    string
	Payload   interface{}
}

type PingPayload struct {
	PingID         int    `json:"ping_id"`
	CurrentRetries int    `json:"current_retries"`
	IsPong         bool   `json:"is_pong"`
	OriginIP       string `json:"origin_ip"`
}

type ChangedFileListPayload struct {
	Filenames []string `json:"filenames"`
}

type FileDeltasPayload struct {
	RevisionID int               `json:"revision_id"`
	Deltas     map[string]string `json:"deltas"`
}

type ConnectionRequestPayload struct {
	OriginIP       string `json:"origin_ip"`
	RepoRemoteName string `json:"repo_remote_name"`
}
