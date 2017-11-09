package test

import (
	"testing"
	//	fs "github.com/elc1798/teatime/fs"
	"fmt"
	tt "github.com/elc1798/teatime"
	diff "github.com/elc1798/teatime/diff"
	"github.com/stretchr/testify/assert"
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
		t.Error("files with same lines marked as modified")
	}
	if !diff.WasModified(file1, file3) {
		t.Error("files with different lines marked as unchanged")
	}

}

func TestDiff(t *testing.T) {
	file1 := tt.File{}
	file2 := tt.File{}
	file1.AppendLine("Lorem ipsum dolor.")
	file2.AppendLine("Lorem dolor sit amet.")
	d := diff.CreateDiff(file1, file2)
	fmt.Println(d)
	assert.Equal(t, file2, diff.ApplyDiff(file1, d), "Diff was incorrect")
}
