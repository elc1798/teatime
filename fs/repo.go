package fs

import (
	"io/ioutil"
	"os"
	"path"
	"strings"

	tt "github.com/elc1798/teatime"
)

type Repo struct {
	Name    string // Name of repository
	RepoDir string // Root directory of repository in Teatime
	RootDir string // Root directory of files being tracked
}

func (this *Repo) GetTrackedDir() string {
	return path.Join(this.RepoDir, tt.TEATIME_TRACKED_DIR)
}

func (this *Repo) GetBackupDir() string {
	return path.Join(this.RepoDir, tt.TEATIME_BACKUP_DIR)
}

func (this *Repo) GetPeerCacheFile() string {
	return path.Join(this.RepoDir, tt.TEATIME_PEER_CACHE)
}

func (this *Repo) Remove() error {
	return os.RemoveAll(this.RepoDir)
}

func InitRepo(name string, pathToRepo string) (*Repo, error) {
	repoDir := path.Join(tt.TEATIME_DEFAULT_HOME, name)

	// Check if Repo exists
	if pathExists(repoDir) {
		return nil, ErrorRepoAlreadyExists(name)
	}

	// Set up directory
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		defer os.RemoveAll(repoDir)

		return nil, err
	}

	// Set up tracked and backup directors
	os.MkdirAll(path.Join(repoDir, tt.TEATIME_TRACKED_DIR), 0755)
	os.MkdirAll(path.Join(repoDir, tt.TEATIME_BACKUP_DIR), 0755)
	os.Create(path.Join(repoDir, tt.TEATIME_PEER_CACHE))

	// Set up DIR_ROOT_STORE (file that stores path to root of repo
	if err := ioutil.WriteFile(path.Join(repoDir, tt.TEATIME_DIR_ROOT_STORE), []byte(pathToRepo), 0644); err != nil {
		defer os.RemoveAll(repoDir)

		return nil, err
	}

	return LoadRepo(name)
}

func LoadRepo(name string) (*Repo, error) {
	repoDir := path.Join(tt.TEATIME_DEFAULT_HOME, name)

	if !pathExists(repoDir) {
		return nil, ErrorNotRepo(name)
	}

	rootDir, err := tt.ReadFile(path.Join(repoDir, tt.TEATIME_DIR_ROOT_STORE))
	if err != nil {
		return nil, err
	}

	return &Repo{
		Name:    name,
		RepoDir: repoDir,
		RootDir: strings.TrimSpace(rootDir[0]),
	}, nil
}

func GetAllRepos() ([]*Repo, error) {
	files, err := ioutil.ReadDir(tt.TEATIME_DEFAULT_HOME)
	if err != nil {
		return nil, err
	}

	repos := make([]*Repo, 0)
	for _, f := range files {
		tmp, err := LoadRepo(f.Name())
		if err != nil {
			return nil, err
		}

		repos = append(repos, tmp)
	}

	return repos, nil
}
