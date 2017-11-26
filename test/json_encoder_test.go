package test

import (
	"reflect"
	"testing"

	encoder "github.com/elc1798/teatime/encode"
)

func testEncodeDecode(s, input interface{}, t *testing.T) interface{} {
	var serializer encoder.Serializer

	serializer, ok := s.(encoder.Serializer)
	if !ok {
		t.Fatalf("Invalid serializer! Type: %v", reflect.TypeOf(s))
	}

	encoded, err := serializer.Serialize(input)
	if err != nil {
		t.Fatalf("Failed to encode! error='%v'", err)
	}

	t.Logf("Serialized: %v", string(encoded))

	// Deserialize and return
	decoded_obj, err := serializer.Deserialize(encoded)
	if err != nil {
		t.Fatalf("Failed to decode! error='%v'", err)
	}
	t.Logf("Decoded: %v", decoded_obj)

	return decoded_obj
}

func TestConnectRequestSerializer(t *testing.T) {
	s1 := encoder.ConnectionRequestSerializer{}
	x1 := encoder.ConnectionRequestPayload("abc")

	decoded_obj := testEncodeDecode(&s1, x1, t)
	decoded := decoded_obj.(encoder.ConnectionRequestPayload)

	if x1 != decoded {
		t.Fatalf("Decoded != Encoded")
	}
}

func TestChangedFileListSerializer(t *testing.T) {
	filenames := []string{
		"testfile1",
		"teatime_is_best_time.lmao",
		"golang_is_better_than_c.lol",
		"cs241 > ece391 lmao",
	}

	s1 := encoder.ChangedFileListSerializer{}
	x1 := encoder.ChangedFileListPayload{
		Filenames: filenames,
	}

	decoded_obj := testEncodeDecode(&s1, x1, t)
	decoded := decoded_obj.(encoder.ChangedFileListPayload).Filenames

	if len(filenames) != len(decoded) {
		t.Fatalf("Decoding error, length mismatch. Expected '%v', got '%v'", len(filenames), len(decoded))
	}

	for i, v := range decoded {
		if v != filenames[i] {
			t.Fatalf("Decoding error at index %v. Expected '%v', got '%v'", i, filenames[i], v)
		}
	}
}

func TestPingSerializer(t *testing.T) {
	s1 := encoder.PingSerializer{}
	x1 := encoder.PingPayload{
		PingID:         0,
		CurrentRetries: 17,
		IsPong:         false,
	}

	decoded_obj := testEncodeDecode(&s1, x1, t)
	decoded := decoded_obj.(encoder.PingPayload)

	if decoded.PingID != 0 || decoded.CurrentRetries != 17 || decoded.IsPong != false {
		t.Fatalf("Decoding error!")
	}
}

func TestFileDeltasSerializer(t *testing.T) {
	s1 := encoder.FileDeltasSerializer{}
	x1 := encoder.FileDeltasPayload{
		RevisionID: 12,
		Deltas: map[string]string{
			"fake_file.txt":   "idk what a diff string looks like",
			"other_file.lmao": "xd kappa",
		},
	}

	decoded_obj := testEncodeDecode(&s1, x1, t)
	decoded := decoded_obj.(encoder.FileDeltasPayload)

	if decoded.RevisionID != 12 || len(decoded.Deltas) != 2 || decoded.Deltas["fake_file.txt"] != "idk what a diff string looks like" ||
		decoded.Deltas["other_file.lmao"] != "xd kappa" {
		t.Fatalf("Decoding error!")
	}
}

func TestInterTeatimeSerializer(t *testing.T) {
	s1 := encoder.InterTeatimeSerializer{}
	x1 := encoder.TeatimeMessage{
		Recipient: "repo1",
		Action:    encoder.ACTION_PING,
		Payload: encoder.PingPayload{
			PingID:         0,
			CurrentRetries: 1,
			IsPong:         true,
		},
	}

	decoded_obj := testEncodeDecode(&s1, x1, t)
	decoded := decoded_obj.(encoder.TeatimeMessage)

	if decoded.Recipient != x1.Recipient || decoded.Action != x1.Action {
		t.Fatalf("Decoding error!")
	}

	y1 := decoded.Payload.(encoder.PingPayload)
	z1 := x1.Payload.(encoder.PingPayload)
	if y1.PingID != z1.PingID || y1.CurrentRetries != z1.CurrentRetries || y1.IsPong != z1.IsPong {
		t.Fatalf("Decoding error!")
	}
}
