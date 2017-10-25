package fs

import (
	"fmt"
	"os"
	"path/filepath"
    "math/rand"
    "errors"
    "diff"
)

const TEATIME_TRACKED_DIR = ".tracked/"

//Might want to use a home directory to track which directories to poll for changes?
const TEATIME_DEFAULT_HOME = "/.teatime/"

/*
 * Creates a hard link to the file at the given path, and places it in the tracked folder
 * for the working directory.  Does not work if not in top level repo directory, though files
 * not in the current directory can be linked.
 *
 * Returns any system error that occurs, or an error signifying that a file with the same
 * name is already being tracked, or that the current directory is not a TeaTime repo.
 */
func AddTrackedFile(path string) error {
	_, file := filepath.Split(path)
	tempLinkName := getTempLinkName(file)
    finalLinkPath := getTrackedFolderPath() + file

    if !pathExists(getTrackedFolderPath()) {
        return ErrorNotRepo()
    }

    if pathExists(finalLinkPath) {
        return ErrorAlreadyTrackingFile(file)
    }

	err := os.Link(path, tempLinkName)
	if err != nil {
		return err
	}

	err = os.Rename(tempLinkName, finalLinkPath)
	return err
}

/*
 * Reads in a file line by line from the given path, and returns a file object (essentially a vector of lines)
 */
func GetFileObjFromFile(path string) (diff.File, error) {
    return diff.File{}, nil
}

/*
 * Returns either the TEATIME_HOME environment variable, or the defaul home if the environment
 * variable is not set.
 */
func getTTHome() string {
    home := os.Getenv("TEATIME_HOME")
    if home == "" {
        return os.Getenv("HOME") + TEATIME_DEFAULT_HOME
    } else {
        return home
    }
}

/*
 * Returns a string for the full path for the tracked folder
 */
func getTrackedFolderPath() string {
	return TEATIME_TRACKED_DIR
}

/*
 * Naive method for generating unique filename for temporary creation/usage.
 */
func getTempLinkName(path string) string {
	return fmt.Sprintf("%s%s%d", path, ".link.", rand.Int31())
}

func ErrorNotRepo() error {
    return errors.New("Working directory is not a TeaTime repo directory.")
}

func ErrorAlreadyTrackingFile(filename string) error {
    return errors.New("Already tracking a file with name: " + filename)
}

func pathExists(path string) bool {
    _, stat_err := os.Stat(path)
    return !os.IsNotExist(stat_err)
}
