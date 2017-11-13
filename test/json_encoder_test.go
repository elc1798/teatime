package test

import (
	"testing"

	p2p "github.com/elc1798/teatime/p2p"
)

func TestChangedFileListSerializer(t *testing.T) {
	filenames := []string{
		"testfile1",
		"teatime_is_best_time.lmao",
		"golang_is_better_than_c.lol",
		"cs241 > ece391 lmao",
	}

	s1 := p2p.ChangedFileListSerializer{}

	// Check if s1 actually inherits p2p.Serializer
	var v1 interface{} = &s1
	if _, ok := v1.(p2p.Serializer); !ok {
		t.Fatalf("ChangedFileSerializer is not valid Serializer!")
	}

	encoded, err := s1.Serialize(filenames)
	if err != nil {
		t.Fatalf("Failed to encode! error='%v'", err)
	}
	t.Logf("Json: %v", string(encoded))

	// Deserialize and check equality
	decoded_obj, err := s1.Deserialize(encoded)
	if err != nil {
		t.Fatalf("Failed to decode! error='%v'", err)
	}

	decoded := decoded_obj.([]string)
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
