package test

import (
	"os"
	"path"
	"testing"
	"time"

	tt "github.com/elc1798/teatime"
)

const pollFILE1 = "file1.txt"
const pollFILE2 = "file2.txt"
const pollFILE3 = "file3.txt"

var pollFILES = []string{pollFILE1, pollFILE2, pollFILE3}

func TestPollForChanges(t *testing.T) {
	// Clear original teatime directory
	tt.ResetTeatime()

	r1, d1, _ := setUpRepos("lmao")
	defer os.RemoveAll(d1)

	fileOrig := tt.File{}
	fileOrig.AppendLine("hello")
	fileOrig.AppendLine("world")

	fileChanged := tt.File{}
	fileChanged.AppendLine("hello")
	fileChanged.AppendLine("universe")

	//Write original files
	for _, f := range pollFILES {
		tt.WriteFileObjToPath(&fileOrig, path.Join(d1, f))
	}

	// Add to tracked
	for _, f := range pollFILES {
		r1.AddFile(f)
	}

	timer1 := time.NewTimer(time.Millisecond * 500)
	changeDetected, resumePolling := r1.StartPollingRepo()
	<-timer1.C

	//Overwrite some of the tracked files
	tt.WriteFileObjToPath(&fileChanged, path.Join(d1, pollFILE3))
	time.Sleep(time.Millisecond * 500)

	select {
	case <-changeDetected:
		r1.WriteBackupFile(pollFILE3)
	default:
		t.Fatal("Failed to detect first set of changes!\n")
	}

	//Overwrite another of the tracked files
	tt.WriteFileObjToPath(&fileChanged, path.Join(d1, pollFILE2))
	time.Sleep(time.Millisecond * 500)

	select {
	case <-changeDetected:
		t.Fatal("Failed to wait until resume signal was sent through channel!\n")
	default:
		resumePolling <- true
	}

	time.Sleep(time.Millisecond * 500)

	select {
	case <-changeDetected:
		return
	default:
		t.Fatal("Failed to resume polling and detect second set of changes!\n")
	}

	r1.Remove()
	os.RemoveAll(tt.TEATIME_DEFAULT_HOME)
}
