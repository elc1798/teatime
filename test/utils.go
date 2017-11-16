package test

import (
	"os"
	"path"

	fs "github.com/elc1798/teatime/fs"
)

var HOME = os.Getenv("HOME")

func setUpRepos(s string) (*fs.Repo, string, error) {
	repoDir := path.Join(HOME, s)
	os.Mkdir(repoDir, 0755)
	r, err := fs.InitRepo(s, repoDir)
	return r, repoDir, err
}
