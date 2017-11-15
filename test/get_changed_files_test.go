package test

import (
	"os"
	"path"
	"testing"

	tt "github.com/elc1798/teatime"
)

const FILE1 = "file1.txt"
const FILE2 = "file2.txt"
const FILE3 = "file3.txt"

var FILES = []string{FILE1, FILE2, FILE3}

func TestGetChangedFiles(t *testing.T) {
	tt.ResetTeatime()

	r, repoDir, err := setUpRepos("test_repo")
	if err != nil {
		t.Fatalf("Error in repo setup: '%v'", err)
	}

	fileOrig := tt.File{}
	fileOrig.AppendLine("hello")
	fileOrig.AppendLine("world")

	fileChanged := tt.File{}
	fileChanged.AppendLine("hello")
	fileChanged.AppendLine("universe")

	//Write original files
	for _, f := range FILES {
		tt.WriteFileObjToPath(&fileOrig, path.Join(repoDir, f))
	}

	// Add to tracked
	for _, f := range FILES {
		r.AddFile(f)
	}

	//Overwrite some of the tracked files
	tt.WriteFileObjToPath(&fileChanged, path.Join(repoDir, FILE1))
	tt.WriteFileObjToPath(&fileChanged, path.Join(repoDir, FILE3))

	changelist, err := r.GetChangedFiles()

	if err != nil {
		t.Fatalf("Failed to get changed files: %v\n", err)
	}

	if len(changelist) != 2 {
		t.Fatalf("Expected 2 changed files.  Actual: %v\n", len(changelist))
	}

	if changelist[0] != changelist[1] {
		if changelist[0] != FILE1 && changelist[0] != FILE3 {
			t.Fatalf("Received unexpected changed file: %v\n", changelist[0])
		}
		if changelist[1] != FILE1 && changelist[1] != FILE3 {
			t.Fatalf("Received unexpected changed file: %v\n", changelist[1])
		}

	} else {
		t.Fatalf("Received repeats in changed file list\n")
	}

	r.Remove()
	os.RemoveAll(tt.TEATIME_DEFAULT_HOME)
	os.RemoveAll(repoDir)
}
