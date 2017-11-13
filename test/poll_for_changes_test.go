package test

import (
	"os"
	"testing"
    "time"

	tt "github.com/elc1798/teatime"
	fs "github.com/elc1798/teatime/fs"
)

const pollFILE1 = "file1.txt"
const pollFILE2 = "file2.txt"
const pollFILE3 = "file3.txt"

func TestPollForChanges(t *testing.T) {
	os.RemoveAll(tt.TEATIME_DEFAULT_HOME)
	os.RemoveAll(tt.TEATIME_TRACKED_DIR)
	os.RemoveAll(tt.TEATIME_BACKUP_DIR)

	defer os.RemoveAll(tt.TEATIME_DEFAULT_HOME)
	defer os.RemoveAll(tt.TEATIME_TRACKED_DIR)
	defer os.RemoveAll(tt.TEATIME_BACKUP_DIR)


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
	tt.WriteFileObjToPath(&fileOrig, tt.TEATIME_TRACKED_DIR+pollFILE1)
	tt.WriteFileObjToPath(&fileOrig, tt.TEATIME_TRACKED_DIR+pollFILE2)
	tt.WriteFileObjToPath(&fileOrig, tt.TEATIME_TRACKED_DIR+pollFILE3)

	//Write backup files to diff with
	err := fs.WriteBackupFile(pollFILE1)
	err = fs.WriteBackupFile(pollFILE2)
	err = fs.WriteBackupFile(pollFILE3)
    if err != nil {
        t.Fatalf("Error backing up: %v\n", err)
    }


	changeDetected, resumePolling := fs.StartPollingRepo(".")
    time.Sleep(fs.POLLING_INTERVAL * 3 * time.Millisecond)

	//Overwrite some of the tracked files
	tt.WriteFileObjToPath(&fileChanged, tt.TEATIME_TRACKED_DIR+pollFILE3)

    time.Sleep(fs.POLLING_INTERVAL * 3 * time.Millisecond)

    select {
    case <-changeDetected:
        fs.WriteBackupFile(pollFILE3)
    default:
        t.Fatal("Failed to detect first set of changes!\n")
    }

	//Overwrite another of the tracked files
    time.Sleep(fs.POLLING_INTERVAL * 2 * time.Millisecond)
	tt.WriteFileObjToPath(&fileChanged, tt.TEATIME_TRACKED_DIR+pollFILE2)

    select {
    case <-changeDetected:
        t.Fatal("Failed to wait until resume signal was sent through channel!\n")
    default:
        resumePolling <- true
    }

    time.Sleep(fs.POLLING_INTERVAL * 3 * time.Millisecond)

    select {
    case <-changeDetected:
        return
    default:
        t.Fatal("Failed to resume polling and detect second set of changes!\n")
    }
}
