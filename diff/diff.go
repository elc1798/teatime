package diff

import "crypto/md5"
import "io"
import "bytes"

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

func CreateDiff(basefile fs.File, newfile fs.File) diff {
    //http://www.xmailserver.org/diff2.pdf
    return diff{}
}

func DiffToSwapfile(d diff) fs.File {
    return fs.File{}
}

func SwapfileToDiff(swapfile fs.File) diff {
    return diff{}
}

//Not sure what to do here...
func HandleMergeConflicts(d diff){
}

