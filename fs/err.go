package fs

import (
	"fmt"
)

func ErrorNotRepo(path string) error {
	return fmt.Errorf("'%v' is not a TeaTime repo.", path)
}

func ErrorAlreadyTrackingFile(filename string) error {
	return fmt.Errorf("File '%v' is already being tracked", filename)
}

func ErrorNotTrackingFile(filename string) error {
	return fmt.Errorf("File '%v' is not currently tracked.", filename)
}

func ErrorRepoAlreadyExists(name string) error {
	return fmt.Errorf("Repo '%v' already exists.", name)
}
