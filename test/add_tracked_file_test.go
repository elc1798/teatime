package test

import (
	"os"
	"testing"

	fs "github.com/elc1798/teatime/fs"
)

const MY_FILE_NAME = "add_tracked_file_test.go"

func TestAddTrackedFile(t *testing.T) {
	os.RemoveAll(fs.TEATIME_DEFAULT_HOME)
	os.RemoveAll(fs.TEATIME_TRACKED_DIR)

	if err := os.Mkdir(fs.TEATIME_DEFAULT_HOME, 0777); err != nil {
		t.Fatalf("Error creating teatime home: %v\n", err)
	}

	if err := os.Mkdir(fs.TEATIME_TRACKED_DIR, 0777); err != nil {
		t.Fatalf("Error creating teatime tracked: %v\n", err)
	}

	if err := fs.AddTrackedFile(MY_FILE_NAME); err != nil {
		t.Fatalf("Failed to add tracked file: %v\n", err)
	}

	os.RemoveAll(fs.TEATIME_DEFAULT_HOME)
	os.RemoveAll(fs.TEATIME_TRACKED_DIR)

}
