package fs

import (
	ioutil "io/ioutil"
	"log"
	"os"
	"path"
	"time"

	tt "github.com/elc1798/teatime"
	diff "github.com/elc1798/teatime/diff"
)

/*
 * Creates a hard link to the file at the given path, and places it in the tracked folder
 * for the working directory.  Does not work if not in top level repo directory, though files
 * not in the current directory can be linked.
 *
 * Returns any system error that occurs, or an error signifying that a file with the same
 * name is already being tracked, or that the current directory is not a TeaTime repo.
 */
func (this *Repo) AddFile(relativePath string) error {
	if !isRepoValid(this) {
		return ErrorNotRepo(this.RepoDir)
	}

	linkName := path.Join(this.GetTrackedDir(), relativePath)
	log.Printf("Created link %v", linkName)

	if pathExists(linkName) {
		return ErrorAlreadyTrackingFile(relativePath)
	}

	err := os.Link(path.Join(this.RootDir, relativePath), linkName)
	if err != nil {
		return err
	}

	defer this.WriteBackupFile(relativePath)
	return nil
}

/*
 * Overwrites the file in the "backup" directory for the given filename from the tracked
 * directory.  Only works if the working directory is a teatime repo.
 *
 * Returns any system error that occurs, or an error signifying that a file with the same
 * name is already being tracked, or that the current directory is not a TeaTime repo.
 */
func (this *Repo) WriteBackupFile(trackedRelPath string) error {
	if !isRepoValid(this) {
		return ErrorNotRepo(this.RepoDir)
	}

	if !pathExists(path.Join(this.GetTrackedDir(), trackedRelPath)) {
		return ErrorNotTrackingFile(trackedRelPath)
	}

	err := CopyFile(
		path.Join(this.GetTrackedDir(), trackedRelPath),
		path.Join(this.GetBackupDir(), trackedRelPath),
	)
	return err
}

/*
 * Returns true if any file in the repo has been changed.
 */
func (this *Repo) haveAnyFilesChanged() (bool, error) {
	//Send out task to check for difference in each file in tracked directory
	files, err := ioutil.ReadDir(this.GetTrackedDir())
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
		go this.fileWasChanged(diffChannels[i], errChannels[i], files[i].Name())
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
func (this *Repo) GetChangedFiles() ([]string, error) {
	//Send out task to check for difference in each file in tracked directory
	files, err := ioutil.ReadDir(this.GetTrackedDir())
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
		go this.fileWasChanged(diffChannels[i], errChannels[i], files[i].Name())
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

func (this *Repo) GetDiffStrings(filesToDiff []string) ([]string, []error) {
	diffStrings := make([]string, len(filesToDiff))
	errors := make([]error, len(filesToDiff))

	for i, fileName := range filesToDiff {
		fileTracked, err := tt.GetFileObjFromFile(path.Join(this.GetTrackedDir(), fileName))
		if err != nil {
			diffStrings[i] = ""
			errors[i] = err
			continue
		}

		fileBackup, err := tt.GetFileObjFromFile(path.Join(this.GetBackupDir(), fileName))
		if err != nil {
			diffStrings[i] = ""
			errors[i] = err
			continue
		}

		diffStrings[i] = diff.CreateDiff(*fileBackup, *fileTracked)
		errors[i] = nil
	}

	return diffStrings, errors
}

/*
 *  Checks for a difference in the backup and tracked file for the given file name, and pushes the results
 *  onto the given channels.

 *  Output:
 *      result : true if the file was modified (backup and tracked file have different hashes), false otherwise
 *      err    : any error that occurred during this process
 */
func (this *Repo) fileWasChanged(result chan bool, errChan chan error, fileName string) {
	fileTracked, err := tt.GetFileObjFromFile(path.Join(this.GetTrackedDir(), fileName))
	if err != nil {
		result <- false
		errChan <- err
		return
	}
	fileBackup, err := tt.GetFileObjFromFile(path.Join(this.GetBackupDir(), fileName))
	if err != nil {
		result <- false
		errChan <- err
		return
	}

	result <- diff.WasModified(*fileBackup, *fileTracked)
	errChan <- err
}

func (this *Repo) getDiffString(result chan string, errChan chan error, fileName string) {
	fileTracked, err := tt.GetFileObjFromFile(path.Join(this.GetTrackedDir(), fileName))
	if err != nil {
		result <- ""
		errChan <- err
		return
	}
	log.Printf("Got tracked file contents")
	fileBackup, err := tt.GetFileObjFromFile(path.Join(this.GetBackupDir(), fileName))
	if err != nil {
		result <- ""
		errChan <- err
		return
	}
	log.Printf("Got backup file contents")

	fdiff := diff.CreateDiff(*fileBackup, *fileTracked)
	log.Printf("Got diff: %v", fdiff)
	result <- fdiff
	errChan <- err
}

/*
 * Pushes true to signalChannel whenever changes are detected, then waits on a value being
 * pushed onto the resumeChannel before resuming polling.
 */
func (this *Repo) pollForChanges(signalChannel chan bool, resumeChannel chan bool) {
	ticker := time.NewTicker(time.Millisecond * 250)
	for {
		changed, _ := this.haveAnyFilesChanged()
		if changed {
			signalChannel <- changed
			<-resumeChannel //Block until resume signal is sent
		}

		<-ticker.C
	}
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
func (this *Repo) StartPollingRepo() (chan bool, chan bool) {
	updateDetectedChannel := make(chan bool)
	resumePollingChannel := make(chan bool)

	go this.pollForChanges(updateDetectedChannel, resumePollingChannel)

	return updateDetectedChannel, resumePollingChannel
}

/*
 * Rewrites the backup file for all of the filenames in the given list in the given repo
 */
func (this *Repo) WriteMultipleBackupFiles(fileNameList []string) []error {
	//Send out task to check for difference in each file in tracked directory
	var numFiles int = len(fileNameList)
	errChannels := make([]chan error, numFiles)
	for i := 0; i < numFiles; i++ {
		errChannels[i] = make(chan error)
		go func(errChan chan error, fileName string) {
			backupErr := this.WriteBackupFile(fileName)
			errChan <- backupErr
		}(errChannels[i], fileNameList[i])
	}

	//Receive the results and build result string array
	errors := make([]error, numFiles)
	for i := 0; i < numFiles; i++ {
		errors[i] = <-errChannels[i]
	}

	return errors
}

/*
 * Patches all files in filesToPatch with the corresponding delta string from diffStrings.
 * Writes all changed to the tracked files.
 *
 * Note: Does not overwrite the backup files
 * Also note: (untested as of now)
 */
func (this *Repo) ApplyDiffs(filesToPatch []string, diffStrings []string) []error {
	//Send out task to apply diff to file
	var numToPatch int = len(filesToPatch)
	errChannels := make([]chan error, numToPatch)
	for i := 0; i < numToPatch; i++ {
		errChannels[i] = make(chan error)
		go func(fileName string, diffString string, errChan chan error) {
			patchErr := this.PatchFile(fileName, diffString)
			errChan <- patchErr
		}(filesToPatch[i], diffStrings[i], errChannels[i])
	}

	//Receive any errors
	errors := make([]error, numToPatch)
	for i := 0; i < numToPatch; i++ {
		errors[i] = <-errChannels[i]
	}

	return errors
}

func (this *Repo) PatchFile(fileName string, diffString string) error {
	filepath := path.Join(this.GetTrackedDir(), fileName)
	basefileptr, err := tt.GetFileObjFromFile(filepath)
	if err != nil {
		return err
	}

	newfileobj := diff.ApplyDiff(*basefileptr, diffString)
	err = tt.WriteFileObjToPath(&newfileobj, filepath)
	return err
}

func (this *Repo) PatchFileMergeConflict(fileName string, conflictingDiffs []string) error {
	if len(conflictingDiffs) < 2 {
		return this.PatchFile(fileName, conflictingDiffs[0])
	}

	filepath := path.Join(this.GetTrackedDir(), fileName)
	basefileptr, err := tt.GetFileObjFromFile(filepath)
	if err != nil {
		return err
	}

	newDiffString := diff.HandleMergeConflicts(*basefileptr, conflictingDiffs[0], conflictingDiffs[1])
	for i := 2; i < len(conflictingDiffs)-1; i++ {
		newDiffString = diff.HandleMergeConflicts(*basefileptr, newDiffString, conflictingDiffs[i])
	}

	newfileobj := diff.ApplyDiff(*basefileptr, newDiffString)
	err = tt.WriteFileObjToPath(&newfileobj, filepath)
	return err

}
