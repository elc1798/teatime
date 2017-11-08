package test

import (
	"os"
	"testing"

	tt "github.com/elc1798/teatime"
	fs "github.com/elc1798/teatime/fs"
)

const GET_CHANGED_FILES_TEST_GO = "get_changed_files_test.go"
const FILE1 = "file1.txt"
const FILE2 = "file2.txt"
const FILE3 = "file3.txt"

func TestGetChangedFiles(t *testing.T) {
	os.RemoveAll(tt.TEATIME_DEFAULT_HOME)
	os.RemoveAll(tt.TEATIME_TRACKED_DIR)
	os.RemoveAll(tt.TEATIME_BACKUP_DIR)

	if err := os.Mkdir(tt.TEATIME_DEFAULT_HOME, 0777); err != nil {
		t.Fatalf("Error creating teatime home: %v\n", err)
	}

	if err := os.Mkdir(tt.TEATIME_TRACKED_DIR, 0777); err != nil {
		t.Fatalf("Error creating teatime tracked: %v\n", err)
	}

	if err := os.Mkdir(tt.TEATIME_BACKUP_DIR, 0777); err != nil {
		t.Fatalf("Error creating teatime backup: %v\n", err)
	}

	fileOrig := tt.File{}
	fileOrig.AppendLine("hello")
	fileOrig.AppendLine("world")

	fileChanged := tt.File{}
	fileChanged.AppendLine("hello")
	fileChanged.AppendLine("universe")

    //Write original files
    tt.WriteFileObjToPath(&fileOrig, tt.TEATIME_TRACKED_DIR + FILE1) 
    tt.WriteFileObjToPath(&fileOrig, tt.TEATIME_TRACKED_DIR + FILE2) 
    tt.WriteFileObjToPath(&fileOrig, tt.TEATIME_TRACKED_DIR + FILE3) 

    //Write backup files to diff with
    fs.WriteBackupFile(FILE1)
    fs.WriteBackupFile(FILE2)
    fs.WriteBackupFile(FILE3)
    fs.WriteBackupFile(GET_CHANGED_FILES_TEST_GO)

    //Overwrite some of the tracked files
    tt.WriteFileObjToPath(&fileChanged, tt.TEATIME_TRACKED_DIR + FILE1) 
    tt.WriteFileObjToPath(&fileOrig, tt.TEATIME_TRACKED_DIR + FILE2) 
    tt.WriteFileObjToPath(&fileChanged, tt.TEATIME_TRACKED_DIR + FILE3) 


    changelist, err := fs.GetChangedFiles(".")

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

	os.RemoveAll(tt.TEATIME_DEFAULT_HOME)
	os.RemoveAll(tt.TEATIME_TRACKED_DIR)
	os.RemoveAll(tt.TEATIME_BACKUP_DIR)
}
