package test

import (
	"os"
	"path"
	"testing"

	tt "github.com/elc1798/teatime"
	fs "github.com/elc1798/teatime/fs"
)

const MY_FILE_NAME = "add_tracked_file_test.go"

func TestAddTrackedFile(t *testing.T) {
	// Clear original teatime directory
	tt.ResetTeatime()

	repoDir := path.Join(os.Getenv("HOME"), "test_repo")
	os.Mkdir(repoDir, 0755)
	if err := fs.InitRepo("test_repo", repoDir); err != nil {
		t.Fatalf("Error creating repo: '%v'", err)
	}

	r, err := fs.LoadRepo("test_repo")
	if err != nil {
		t.Fatalf("Error loading repo: '%v'", err)
	}

	if r.Name != "test_repo" || r.RepoDir != path.Join(tt.TEATIME_DEFAULT_HOME, "test_repo") || r.RootDir != repoDir {
		t.Fatalf("Invalid data: '%v'", r)
	}

	r.Remove()
	os.RemoveAll(tt.TEATIME_DEFAULT_HOME)
	os.RemoveAll(repoDir)
}
