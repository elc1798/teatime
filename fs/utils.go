package fs

import (
	"fmt"
	"io"
	"math/rand"
	"os"
)

func getTempLinkName(path string) string {
	return fmt.Sprintf("%s%s%d", path, ".link.", rand.Int31())
}

func pathExists(path string) bool {
	_, stat_err := os.Stat(path)
	return !os.IsNotExist(stat_err)
}

func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	return nil
}

func isRepoValid(r *Repo) bool {
	return pathExists(r.GetTrackedDir()) && pathExists(r.GetBackupDir())
}
