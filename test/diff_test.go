package diff

import (
    "testing"
	fs "github.com/elc1798/teatime/fs"
	dmp "github.com/sergi/go-diff/diffmatchpatch"
    "github.com/stretchr/testify/assert"
)

func TestWasModified(t *testing.T){
    file1 := fs.File{}
    file1.AppendLine("hello")
    file1.AppendLine("world")

    file2 := fs.File{}
    file2.AppendLine("hello")
    file2.AppendLine("world")

    file3 := fs.File{}
    file3.AppendLine("goodbye")
    file3.AppendLine("world")

    if WasModified(file1, file2) {
        t.Error("files with same lines marked as modified")
    }
    if !WasModified(file1, file3) {
        t.Error("files with different lines marked as unchanged")
    }

}

func TestDiff(t *testing.T){
    file1 := fs.File{}
    file2 := fs.File{}
    //file1.AppendLine("Lorem ipsum dolor.")
    //file2.AppendLine("Lorem dolor sit amet.")
    file1.AppendLine("Lorem ipsum")
    file2.AppendLine("Lorem dolor")
    diffs1 := []dmp.Diff { dmp.Diff {Type:0, Text:"Lorem " }, dmp.Diff {Type:-1, Text:"ipsum"}, dmp.Diff {Type:1, Text:"dolor"} }

    assert.Equal(t, CreateDiff(file1, file2), diffs1, "Diff was incorrect")
}

