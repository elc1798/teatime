package test

import (
	"testing"

    tt   "github.com/elc1798/teatime"
	diff "github.com/elc1798/teatime/diff"
)

func TestWasModified(t *testing.T) {
	file1 := tt.File{}
	file1.AppendLine("hello")
	file1.AppendLine("world")

	file2 := tt.File{}
	file2.AppendLine("hello")
	file2.AppendLine("world")

	file3 := tt.File{}
	file3.AppendLine("goodbye")
	file3.AppendLine("world")

	if diff.WasModified(file1, file2) {
		t.Fatal("files with same lines marked as modified")
	}
	if !diff.WasModified(file1, file3) {
		t.Fatal("files with different lines marked as unchanged")
	}

}
