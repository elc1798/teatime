package fs

import (
	"errors"
	"fmt"
	"io"
	ioutil "io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	tt "github.com/elc1798/teatime"
	diff "github.com/elc1798/teatime/diff"
)

//TODO: Applying set of diffs to files / backing up all tracked files

const POLLING_INTERVAL = 100

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
 * Start polling a repo for changes.  Returns two bool channels.
 *
 * The first returned channel will have true pushed to it whenever a change is
 * detected in the repo.
 *
 * After handling the detected changes, push a value to the second returned channel
 * to resume polling.
 */
func StartPollingRepo(pathToRepo string) (chan bool, chan bool) {
	updateDetectedChannel := make(chan bool)
	resumePollingChannel := make(chan bool)

	go pollForChanges(pathToRepo, updateDetectedChannel, resumePollingChannel)

	return updateDetectedChannel, resumePollingChannel
}

/*
 * Returns true if any file in the repo has been changed.
 */
func haveAnyFilesChanged(pathToRepo string) (bool, error) {
	err := os.Chdir(pathToRepo)
	if err != nil {
		return false, err
	}

	//Send out task to check for difference in each file in tracked directory
	files, err := ioutil.ReadDir(tt.TEATIME_TRACKED_DIR)
	var numTracked int = len(files)
	diffChannels := make([]chan bool, numTracked)
	errChannels := make([]chan error, numTracked)
	for i := 0; i < numTracked; i++ {
		diffChannels[i] = make(chan bool)
		errChannels[i] = make(chan error)
		if files[i].IsDir() {
			diffChannels[i] <- false
			continue
		}
		go fileWasChanged(diffChannels[i], errChannels[i], files[i].Name())
	}

	//Receive the results
	for i := 0; i < numTracked; i++ {
		wasChanged := <-diffChannels[i]

		newErr := <-errChannels[i]
		if newErr != nil && err == nil {
			err = newErr
		}

		if wasChanged {
			return true, err
		}
	}

	return false, err
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
	files, err := ioutil.ReadDir(tt.TEATIME_TRACKED_DIR)
	var numTracked int = len(files)
	diffChannels := make([]chan bool, numTracked)
	errChannels := make([]chan error, numTracked)
	for i := 0; i < numTracked; i++ {
		diffChannels[i] = make(chan bool)
		errChannels[i] = make(chan error)
		if files[i].IsDir() {
			diffChannels[i] <- false
			continue
		}
		go fileWasChanged(diffChannels[i], errChannels[i], files[i].Name())
	}

	//Receive the results and build result string array
	var changedFiles []string
	for i := 0; i < numTracked; i++ {
		wasChanged := <-diffChannels[i]
		if wasChanged {
			changedFiles = append(changedFiles, files[i].Name())
		}
		newErr := <-errChannels[i]
		if newErr != nil && err == nil {
			err = newErr
		}
	}

	return changedFiles, err
}

func GetDiffStrings(pathToRepo string, filesToDiff []string) ([]string, error) {
	err := os.Chdir(pathToRepo)
	if err != nil {
		return nil, err
	}

	//Send out task to build diffstring from file
	var numToDiff int = len(filesToDiff)
	diffChannels := make([]chan string, numToDiff)
	errChannels := make([]chan error, numToDiff)
	for i := 0; i < numToDiff; i++ {
		diffChannels[i] = make(chan string)
		errChannels[i] = make(chan error)
		go getDiffString(diffChannels[i], errChannels[i], filesToDiff[i])
	}

	//Receive the results and build result string array
	var diffStrings []string
	for i := 0; i < numToDiff; i++ {
		newErr := <-errChannels[i]
		if newErr != nil && err == nil {
			err = newErr
		}
		diffValue := <-diffChannels[i]
		diffStrings = append(diffStrings, diffValue)
	}

	return diffStrings, err
}

/*
 * Pushes true to signalChannel whenever changes are detected
 */
func pollForChanges(pathToRepo string, signalChannel chan bool, resumeChannel chan bool) {
	for {
		changed, _ := haveAnyFilesChanged(pathToRepo)
		if changed {
			signalChannel <- changed
			<-resumeChannel //Block until resume signal is sent
			continue
		}
		time.Sleep(POLLING_INTERVAL * time.Millisecond)
	}
}

/*
 *  Checks for a difference in the backup and tracked file for the given file name, and pushes the results
 *  onto the given channels.

 *  Output:
 *      result : true if the file was modified (backup and tracked file have different hashes), false otherwise
 *      err    : any error that occurred during this process
 */
func fileWasChanged(result chan bool, errChan chan error, fileName string) {
	fileTracked, err := tt.GetFileObjFromFile(tt.TEATIME_TRACKED_DIR + fileName)
	if err != nil {
		result <- false
		errChan <- err
		return
	}
	fileBackup, err := tt.GetFileObjFromFile(tt.TEATIME_BACKUP_DIR + fileName)
	if err != nil {
		result <- false
		errChan <- err
		return
	}

	result <- diff.WasModified(*fileBackup, *fileTracked)
	errChan <- err
}

func getDiffString(result chan string, errChan chan error, fileName string) {
	fileTracked, err := tt.GetFileObjFromFile(tt.TEATIME_TRACKED_DIR + fileName)
	if err != nil {
		result <- ""
		errChan <- err
		return
	}
	fileBackup, err := tt.GetFileObjFromFile(tt.TEATIME_BACKUP_DIR + fileName)
	if err != nil {
		result <- ""
		errChan <- err
		return
	}

	result <- diff.CreateDiff(*fileBackup, *fileTracked)
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
