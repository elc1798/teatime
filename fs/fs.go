package fs

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
    "io"
    "bufio"
)

const TEATIME_TRACKED_DIR = ".tracked/"
const TEATIME_BACKUP_DIR = ".backup/"

//Might want to use a home directory to track which directories to poll for changes?
const TEATIME_DEFAULT_HOME = "/.teatime/"

/*
 * File object struct definition.  Used for diffs.
 */
type File struct {
    lineSlice []string
}

func (f *File) GetLine(i int) string {
    return f.lineSlice[i]
}

func (f *File) SetLine(i int, s string) {
    f.lineSlice[i] = s
}

func (f *File) AppendLine(s string) {
    f.lineSlice = append(f.lineSlice, s)
}

func (f *File) NumLines() int {
    return len(f.lineSlice)
}

/*
 * Reads in a file line by line from the given path, and returns a file object 
 * (essentially a vector of lines)
 */
func GetFileObjFromFile(path string) (*File, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    fileObjPtr := new(File)
    for scanner.Scan() {
        fileObjPtr.AppendLine(scanner.Text())
    }

    err = scanner.Err()
    if err != nil {
        return nil, err
    }

	return fileObjPtr, nil
}

func WriteFileObjToPath(fileObj *File, path string) error {
    file, err := os.Create(path)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := bufio.NewWriter(file)
    for i := 0; i < fileObj.NumLines(); i++ {
        _, err = writer.WriteString( fileObj.GetLine(i) + "\n")
        if err != nil {
            return err
        }
    }
    writer.Flush()
    return err
}

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

	err := CopyFile(getTrackedFolderPath() + trackedFileName, getBackupFolderPath() + trackedFileName)
    return err
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
 * Returns a string for the full path for the tracked folder
 */
func getBackupFolderPath() string {
	return TEATIME_BACKUP_DIR
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
