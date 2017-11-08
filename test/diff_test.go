package test

import (
	"testing"

	diff "github.com/elc1798/teatime/diff"
	fs "github.com/elc1798/teatime/fs"
)

func TestWasModified(t *testing.T) {
	file1 := fs.File{}
	file1.AppendLine("hello")
	file1.AppendLine("world")

	file2 := fs.File{}
	file2.AppendLine("hello")
	file2.AppendLine("world")

	file3 := fs.File{}
	file3.AppendLine("goodbye")
	file3.AppendLine("world")

	if diff.WasModified(file1, file2) {
		t.Fatal("files with same lines marked as modified")
	}
	if !diff.WasModified(file1, file3) {
		t.Fatal("files with different lines marked as unchanged")
	}

}
