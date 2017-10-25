package diff

import "crypto/md5"
import "io"
import "bytes"

type diff struct {
}

type File struct {
    lineSlice []string
}

func (f *File) GetLine(i int) string {
    return f.lineSlice[i]
}

func (f *File) SetLine(i int, s string) {
    f.lineSlice[i] = s
}

func (f *File) AppendLine(s string) {
    f.lineSlice = append(f.lineSlice, s)
}

func (f *File) NumLines() int {
    return len(f.lineSlice)
}

func fileHash(f File) []byte {
    h := md5.New()
    for i := 0; i < f.NumLines(); i++ {
        io.WriteString(h, f.GetLine(i))
    }
    return h.Sum(nil)
}

func WasModified(basefile File, newfile File) bool {
    return !bytes.Equal(fileHash(basefile), fileHash(newfile))
}

func CreateDiff(basefile File, newfile File) diff {
    //http://www.xmailserver.org/diff2.pdf
    return diff{}
}

func DiffToSwapfile(d diff) File {
    return File{}
}

func SwapfileToDiff(swapfile File) diff {
    return diff{}
}

//Not sure what to do here...
func HandleMergeConflicts(d diff){
}

