package diff

import (
	"bytes"
	"crypto/md5"
	"io"

	tt "github.com/elc1798/teatime/"
)

type diff struct {
}

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

func CreateDiff(basefile tt.File, newfile tt.File) diff {
	//http://www.xmailserver.org/diff2.pdf
	return diff{}
}

func DiffToSwapfile(d diff) tt.File {
	return tt.File{}
}

func SwapfileToDiff(swapfile tt.File) diff {
	return diff{}
}

//Not sure what to do here...
func HandleMergeConflicts(d diff) {
}
