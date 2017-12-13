package diff

import (
	"bytes"
	"crypto/md5"
	//"fmt"
	tt "github.com/elc1798/teatime"
	dmp "github.com/sergi/go-diff/diffmatchpatch"
	"io"
	"strings"
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

func chooseDiff(diff1 dmp.Diff, diff2 dmp.Diff) dmp.Diff {
	if diff1.Type == 0 {
		return diff2
	}
	return diff1
}
func HandleMergeConflicts(basefile tt.File, delta1 string, delta2 string) string {
	d := dmp.New()
	var newDiffs []dmp.Diff

	diffs2, _ := d.DiffFromDelta(basefile.ToString(), delta1)
	diffs1, _ := d.DiffFromDelta(basefile.ToString(), delta2)
	if delta1 > delta2 {
		diffs1, _ = d.DiffFromDelta(basefile.ToString(), delta1)
		diffs2, _ = d.DiffFromDelta(basefile.ToString(), delta2)
	}

	for index1, index2 := 0, 0; index1 < len(diffs1) || index2 < len(diffs2); {
		diff1 := diffs1[index1]
		diff2 := diffs2[index2]
		if len(diffs1) == index1 {
			newDiffs = append(newDiffs, diff2)
			index2++
			continue
		}
		if len(diffs2) == index2 {
			newDiffs = append(newDiffs, diff1)
			index1++
			continue
		}
		if diff1.Text == diff2.Text {
			newDiffs = append(newDiffs, chooseDiff(diff1, diff2))
			index1++
			index2++
			continue
		}
		if strings.HasSuffix(diff1.Text, diff2.Text) && diff1.Type != 1 {
			newDiffs = append(newDiffs, diff2)
			index1++
			index2++
			continue
		}
		if strings.HasSuffix(diff2.Text, diff1.Text) && diff2.Type != 1 {
			newDiffs = append(newDiffs, diff1)
			index1++
			index2++
			continue
		}
		if strings.Contains(diff1.Text, diff2.Text) && diff1.Type != 1 {
			newDiffs = append(newDiffs, diff2)
			index2++
			continue
		}
		if strings.Contains(diff2.Text, diff1.Text) && diff2.Type != 1 {
			newDiffs = append(newDiffs, diff1)
			index1++
			continue
		}

		if strings.HasSuffix(diff1.Text, diff2.Text) && diff1.Type == 1 {
			newDiffs = append(newDiffs, diff1)
			index1++
			index2++
			continue
		}
		if strings.HasSuffix(diff2.Text, diff1.Text) && diff2.Type == 1 {
			newDiffs = append(newDiffs, diff2)
			index1++
			index2++
			continue
		}
		if strings.Contains(diff1.Text, diff2.Text) && diff1.Type == 1 {
			index2++
			continue
		}
		if strings.Contains(diff2.Text, diff1.Text) && diff2.Type == 1 {
			index1++
			continue
		}
		newDiffs = append(newDiffs, diff1)
		newDiffs = append(newDiffs, diff2)
		index1++
		index2++
	}
	newfileString := d.DiffText2(newDiffs)
	mergediffs := d.DiffMain(basefile.ToString(), newfileString, false)
	return d.DiffToDelta(mergediffs)
}
