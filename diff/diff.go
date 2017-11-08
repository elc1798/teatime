package diff

import (
	"bytes"
	"crypto/md5"
	"io"
    "fmt"
    "github.com/sergi/go-diff/diffmatchpatch"
	fs "github.com/elc1798/teatime/fs"
)

type diff struct {
}

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

func CreateDiff(basefile file, newfile file) diff {
    //http://www.xmailserver.org/diff2.pdf
    text1 := "Lorem ipsum dolor."
    text2 := "Lorem dolor sit amet."
    dmp := diffmatchpatch.New()
    diffs := dmp.DiffMain(text1, text2, false);
    fmt.Println(dmp.DiffPrettyText(diffs));

    return diff{}
}

func DiffToSwapfile(d diff) fs.File {
	return fs.File{}
}

func SwapfileToDiff(swapfile fs.File) diff {
	return diff{}
}

//Not sure what to do here...
func HandleMergeConflicts(d diff) {
}
