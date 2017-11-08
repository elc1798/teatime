package diff

import (
	"bytes"
	"crypto/md5"
	//"fmt"
	fs "github.com/elc1798/teatime/fs"
	dmp "github.com/sergi/go-diff/diffmatchpatch"
	"io"
)

func fileHash(f fs.File) []byte {
	h := md5.New()
	for i := 0; i < f.NumLines(); i++ {
		io.WriteString(h, f.GetLine(i))
	}
	return h.Sum(nil)
}

func WasModified(basefile fs.File, newfile fs.File) bool {
	return !bytes.Equal(fileHash(basefile), fileHash(newfile))
}

func CreateDiff(basefile fs.File, newfile fs.File) []dmp.Diff {
	//http://www.xmailserver.org/diff2.pdf
	d := dmp.New()
	diffs := d.DiffMain(basefile.ToString(), newfile.ToString(), false)
	//fmt.Println(dmp.DiffPrettyText(diffs))
	return diffs
}

func DiffToSwapfile(d []dmp.Diff) fs.File {
	return fs.File{}
}

func SwapfileToDiff(swapfile fs.File) []dmp.Diff {
    //d := []dmp.diff
    var d []dmp.Diff
    return d
}

//Not sure what to do here...
func HandleMergeConflicts(d []dmp.Diff) {
}
