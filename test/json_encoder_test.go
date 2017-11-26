package test

import (
	"testing"

	encoder "github.com/elc1798/teatime/encode"
)

func TestChangedFileListSerializer(t *testing.T) {
	filenames := []string{
		"testfile1",
		"teatime_is_best_time.lmao",
		"golang_is_better_than_c.lol",
		"cs241 > ece391 lmao",
	}

	s1 := encoder.ChangedFileListSerializer{}

	// Check if s1 actually inherits Serializer
	var v1 interface{} = &s1
	if _, ok := v1.(encoder.Serializer); !ok {
		t.Fatalf("ChangedFileSerializer is not valid Serializer!")
	}

	encoded, err := s1.Serialize(encoder.ChangedFileListPayload{
		Filenames: filenames,
	})
	if err != nil {
		t.Fatalf("Failed to encode! error='%v'", err)
	}
	t.Logf("Json: %v", string(encoded))

	// Deserialize and check equality
	decoded_obj, err := s1.Deserialize(encoded)
	if err != nil {
		t.Fatalf("Failed to decode! error='%v'", err)
	}

	decoded := decoded_obj.(encoder.ChangedFileListPayload).Filenames
	t.Logf("Decoded: %v", decoded)
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

	// Check if s1 actually inherits Serializer
	var v1 interface{} = &s1
	if _, ok := v1.(encoder.Serializer); !ok {
		t.Fatalf("PingSerializer is not valid Serializer!")
	}

	encoded, err := s1.Serialize(encoder.PingPayload{
		PingID:         0,
		CurrentRetries: 17,
		IsPong:         false,
	})
	if err != nil {
		t.Fatalf("Failed to encode! error='%v'", err)
	}
	t.Logf("Json: %v", string(encoded))

	// Deserialize and check equality
	decoded_obj, err := s1.Deserialize(encoded)
	if err != nil {
		t.Fatalf("Failed to decode! error='%v'", err)
	}

	decoded := decoded_obj.(encoder.PingPayload)
	t.Logf("Decoded: %v", decoded)

	if decoded.PingID != 0 || decoded.CurrentRetries != 17 || decoded.IsPong != false {
		t.Fatalf("Decoding error!")
	}
}

func TestFileDeltasSerializer(t *testing.T) {
	s1 := encoder.FileDeltasSerializer{}

	// Check if s1 actually inherits Serializer
	var v1 interface{} = &s1
	if _, ok := v1.(encoder.Serializer); !ok {
		t.Fatalf("FileDeltasSerializer is not valid Serializer!")
	}

	encoded, err := s1.Serialize(encoder.FileDeltasPayload{
		RevisionID: 12,
		Deltas: map[string]string{
			"fake_file.txt":   "idk what a diff string looks like",
			"other_file.lmao": "xd kappa",
		},
	})
	if err != nil {
		t.Fatalf("Failed to encode! error='%v'", err)
	}
	t.Logf("Json: %v", string(encoded))

	// Deserialize and check equality
	decoded_obj, err := s1.Deserialize(encoded)
	if err != nil {
		t.Fatalf("Failed to decode! error='%v'", err)
	}

	decoded := decoded_obj.(encoder.FileDeltasPayload)
	t.Logf("Decoded: %v", decoded)

	if decoded.RevisionID != 12 || len(decoded.Deltas) != 2 || decoded.Deltas["fake_file.txt"] != "idk what a diff string looks like" ||
		decoded.Deltas["other_file.lmao"] != "xd kappa" {
		t.Fatalf("Decoding error!")
	}
}
