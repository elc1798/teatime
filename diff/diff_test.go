package diff

import "testing"

func TestWasModified(t *testing.T){
    file1 := file{}
    file1.AppendLine("hello")
    file1.AppendLine("world")

    file2 := file{}
    file2.AppendLine("hello")
    file2.AppendLine("world")

    file3 := file{}
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
    file1 := file{}
    CreateDiff(file1, file1)
}

