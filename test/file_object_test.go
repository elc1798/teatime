package test

import (
	"bufio"
	"crypto/rand"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"strings"
	"testing"

	fs "github.com/elc1798/teatime/fs"
)

const TEST_FILE_NAME = "file_object_test.tmp"

func setUpTestFile() string {
	os.Remove(TEST_FILE_NAME)
	os.Remove(TEST_FILE_NAME + ".out")
	i, err := rand.Int(rand.Reader, big.NewInt(10000000))

	randString := i.String()
	err = ioutil.WriteFile(TEST_FILE_NAME, []byte(randString), 0644)
	if err != nil {
		log.Fatalf("Failed to write test file: %v\n", err)
	}

	return randString
}

func readFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lines := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, nil
}

func TestFileObjectWriting(t *testing.T) {
	generatedString := setUpTestFile()

	fileObj, err := fs.GetFileObjFromFile(TEST_FILE_NAME)
	if err != nil {
		t.Fatalf("Error in GetFileObjFromFile: %v\n", err)
	}

	err = fs.WriteFileObjToPath(fileObj, TEST_FILE_NAME+".out")
	if err != nil {
		t.Fatalf("Error in WriteFileObjToPath: %v\n", err)
	}

	// Check if new file contains same contents
	lines, err := readFile(TEST_FILE_NAME + ".out")
	if len(lines) != 1 {
		t.Fatalf("Invalid contents: %v\n", strings.Join(lines, "\n"))
	}

	if lines[0] != generatedString {
		t.Fatalf("Invalid contents: %v\n", strings.Join(lines, "\n"))
	}

	os.Remove(TEST_FILE_NAME)
	os.Remove(TEST_FILE_NAME + ".out")
}
