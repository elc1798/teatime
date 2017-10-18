package diff
/*
import (
    "io/ioutil"
)
*/

type diff struct {
}

type file struct {
}

func WasModified(basefile file, newfile file) bool {
    return false;
}

func CreateDiff(basefile file, newfile file) diff {
    //http://www.xmailserver.org/diff2.pdf
    return diff{};
}

func DiffToSwapfile(d diff) file {
    return file{};
}

func SwapfileToDiff(swapfile file) diff {
    return diff{};
}

//Not sure what to do here...
func HandleMergeConflicts(d diff){
}

