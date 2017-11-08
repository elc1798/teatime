package fs

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"

    diff "github.com/elc1798/teatime/diff"
	tt "github.com/elc1798/teatime"
)



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
 * Overwrites the file in the "backup" directory for the given filename from the tracked
 * directory.  Only works if the working directory is a teatime repo.
 *
 * Returns any system error that occurs, or an error signifying that a file with the same
 * name is already being tracked, or that the current directory is not a TeaTime repo.
 */
func WriteBackupFile(trackedFileName string) error {
	if !pathExists(getTrackedFolderPath()) || !pathExists(getBackupFolderPath()) {
		return ErrorNotRepo()
	}

	if !pathExists(getTrackedFolderPath() + trackedFileName) {
		return ErrorNotTrackingFile(trackedFileName)
	}

	err := CopyFile(getTrackedFolderPath()+trackedFileName, getBackupFolderPath()+trackedFileName)
	return err
}

/*
 * Returns a list of filenames for all the files that have changes in the repo directory given.
 * Uses goroutines to work in parallel for each file in the directory.
 */
func GetChangedFiles(pathToRepo string) ([]string, error) {
    err := os.Chdir(pathToRepo)
    if err != nil {
        return nil, err
    }

    //Send out task to check for difference in each file in tracked directory
    files, err := io.ReadDir(tt.TEATIME_TRACKED_DIR)
    var numTracked int = len(files)
    var diffChannels [numTracked]chan bool 
    var errChannels [numTracked]chan error
    for i := 0; i < numTracked; i++ {
        diffChannels[i] = make(chan bool)
        errChannels[i] = make(chan error)
        if files[i].IsDir() {
            diffChannels[i] <- false
            continue
        }
        go FileWasChanged(diffChannels[i], errChannels[i], files[i].Name())
    }

    //Receive the results and build result string array
    var changedFiles []string
    for i := 0; i < numTracked; i++ {
        wasChanged := <-diffChannels[i]
        if wasChanged {
            changedFiles = append(changedFiles, files[i].Name()) 
        }
        newErr := <-errChannels[i]
        if (newErr != nil && err == nil) {
            err = newErr
        }
    }


    return changedFiles, err
}

/*
 *  Checks for a difference in the backup and tracked file for the given file name, and pushes the results
 *  onto the given channels.  

 *  Output:
 *      result : true if the file was modified (backup and tracked file have different hashes), false otherwise
 *      err    : any error that occurred during this process
 */
func FileWasChanged(result chan bool, errChan chan error, fileName string) {
    fileTracked, err := tt.GetFileObjFromFile( tt.TEATIME_TRACKED_DIR + fileName )
    if err != nil {
        result <- false
        errChan <- err
    }
    fileBackup, err := tt.GetFileObjFromFile( tt.TEATIME_BACKUP_DIR + fileName )
    if err != nil {
        result <- false
        errChan <- err
    }

    result <- diff.WasModified(fileBackup, fileTracked)
    errChan <- err
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

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return out.Close()
}

/*
 * Returns either the TEATIME_HOME environment variable, or the default home if the environment
 * variable is not set.
 */
func getTTHome() string {
	home := os.Getenv("TEATIME_HOME")
	if home == "" {
		return tt.TEATIME_DEFAULT_HOME
	} else {
		return home
	}
}

/*
 * Returns a string for the full path for the tracked folder
 */
func getTrackedFolderPath() string {
	return tt.TEATIME_TRACKED_DIR
}

/*
 * Returns a string for the full path for the tracked folder
 */
func getBackupFolderPath() string {
	return tt.TEATIME_BACKUP_DIR
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

func ErrorNotTrackingFile(filename string) error {
	return errors.New("Not tracking a file with name: " + filename)
}

func pathExists(path string) bool {
	_, stat_err := os.Stat(path)
	return !os.IsNotExist(stat_err)
}
