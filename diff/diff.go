package diff

import "crypto/md5"
import "io"
import "bytes"

type diff struct {
}

type file struct {
    lineSlice []string
}

func (f *file) GetLine(i int) string {
    return f.lineSlice[i]
}

func (f *file) SetLine(i int, s string) {
    f.lineSlice[i] = s
}

func (f *file) AppendLine(s string) {
    f.lineSlice = append(f.lineSlice, s)
}

func (f *file) NumLines() int {
    return len(f.lineSlice)
}

func fileHash(f file) []byte {
    h := md5.New()
    for i := 0; i < f.NumLines(); i++ {
        io.WriteString(h, f.GetLine(i))
    }
    return h.Sum(nil)
}

func WasModified(basefile file, newfile file) bool {
    return !bytes.Equal(fileHash(basefile), fileHash(newfile))
}

func CreateDiff(basefile file, newfile file) diff {
    //http://www.xmailserver.org/diff2.pdf
    return diff{}
}

func DiffToSwapfile(d diff) file {
    return file{}
}

func SwapfileToDiff(swapfile file) diff {
    return diff{}
}

//Not sure what to do here...
func HandleMergeConflicts(d diff){
}

