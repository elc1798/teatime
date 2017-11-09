package diff

import (
	"bytes"
	"crypto/md5"
	//"fmt"
	tt "github.com/elc1798/teatime"
	dmp "github.com/sergi/go-diff/diffmatchpatch"
	"io"
)

func fileHash(f tt.File) []byte {
	h := md5.New()
	for i := 0; i < f.NumLines(); i++ {
		io.WriteString(h, f.GetLine(i))
	}
	return h.Sum(nil)
}

func WasModified(basefile tt.File, newfile tt.File) bool {
	return !bytes.Equal(fileHash(basefile), fileHash(newfile))
}

func CreateDiff(basefile tt.File, newfile tt.File) string {
	//http://www.xmailserver.org/diff2.pdf
	//https://godoc.org/github.com/sergi/go-diff/diffmatchpatch
	d := dmp.New()
	diffs := d.DiffMain(basefile.ToString(), newfile.ToString(), false)
	return d.DiffToDelta(diffs)
}

func ApplyDiff(basefile tt.File, delta string) tt.File {
	d := dmp.New()
	diffs, _ := d.DiffFromDelta(basefile.ToString(), delta)
	newfileString := d.DiffText2(diffs)
	newfile := tt.File{}
	newfile.FromString(newfileString)
	return newfile
}

//Not sure what to do here...
func HandleMergeConflicts(basefile tt.File, d1 []dmp.Diff, d2 []dmp.Diff) tt.File {
	return tt.File{}
}
