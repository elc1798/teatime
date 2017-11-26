package encode

import (
	json "encoding/json"
	"errors"
)

const ACTION_PING = "ping"
const ACTION_FILE_LIST = "changed_f_list"
const ACTION_DELTAS = "deltas"

type routedPackage struct {
	Destination string `json:"destination"`
	Action      string `json:"action"`
	Data        string `json:"data"`
}

type InterTeatimeSerializer struct{}

/*
 * v must be a map[string]string with "Destination", "Action", and "Data" as
 * keys. If any keys are missing, a "Bad Payload" error is returned. If v is not
 * of type map[string]string, a "Invalid Payload Type" error is returned.
 */
func (this *InterTeatimeSerializer) Serialize(v interface{}) ([]byte, error) {
	dict, ok := v.(map[string]string)
	if !ok {
		return nil, errors.New("Invalid Payload Type")
	}

	checkKey := func(d map[string]string, k string) bool {
		_, ok := d[k]
		return ok
	}

	for _, v := range []string{"Destination", "Action", "Data"} {
		if !checkKey(dict, v) {
			return nil, errors.New("Bad Payload")
		}
	}

	return json.Marshal(routedPackage{
		Destination: dict["Destination"],
		Action:      dict["Action"],
		Data:        dict["Data"],
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
	var payloadSerializer Serializer

	switch pack.Action {
	case ACTION_PING:
		payloadSerializer = &PingSerializer{}
		break
	case ACTION_FILE_LIST:
		payloadSerializer = &ChangedFileListSerializer{}
		break
	case ACTION_DELTAS:
		payloadSerializer = &FileDeltasSerializer{}
		break
	default:
		return nil, errors.New("Invalid action received!")
	}

	payload, err := payloadSerializer.Deserialize([]byte(pack.Data))
	if err != nil {
		return nil, err
	}

	// Cast the interface{} to desired struct for convenience. If any of these
	// errors, our program shouldn't be running anyway, so no errors are caught
	switch pack.Action {
	case ACTION_PING:
		newMessage.Payload = payload.(PingPayload)
		break
	case ACTION_FILE_LIST:
		newMessage.Payload = payload.(ChangedFileListPayload)
		break
	case ACTION_DELTAS:
		newMessage.Payload = payload.(FileDeltasPayload)
		break
	default:
		return nil, errors.New("Invalid action received!")
	}
	return newMessage, nil
}

type TeatimeMessage struct {
	Recipient string
	Action    string
	Payload   interface{}
}

type PingPayload struct {
	PingID         int  `json:"ping_id"`
	CurrentRetries int  `json:"current_retries"`
	IsPong         bool `json:"is_pong"`
}

type ChangedFileListPayload struct {
	Filenames []string `json:"filenames"`
}

type FileDeltasPayload struct {
	RevisionID int               `json:"revision_id"`
	Deltas     map[string]string `json:"deltas"`
}
