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

func CreateDiff(basefile tt.File, newfile tt.File) []dmp.Diff {
	//http://www.xmailserver.org/diff2.pdf
	d := dmp.New()
	diffs := d.DiffMain(basefile.ToString(), newfile.ToString(), false)
	//fmt.Println(dmp.DiffPrettyText(diffs))
	return diffs
}

func SwapfileToDiff(swapfile tt.File) []dmp.Diff {
	//d := []dmp.diff
	var d []dmp.Diff
	return d
}

func DiffToSwapfile(d dmp.Diff) tt.File {
	return tt.File{}
}

//Not sure what to do here...
func HandleMergeConflicts(d []dmp.Diff) {
}
