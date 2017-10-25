package diff

//https://stackoverflow.com/questions/1726336/how-do-i-use-a-generic-vector-in-go

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

func WasModified(basefile file, newfile file) bool {
    return false
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

