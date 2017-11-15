package teatime

import (
	"bufio"
	"bytes"
	"os"
	"strings"
)

// Project constants

const TEATIME_TRACKED_DIR = ".tracked/"
const TEATIME_BACKUP_DIR = ".backup/"
const TEATIME_PEER_CACHE = "/peer_cache"

var TEATIME_DEFAULT_HOME = os.Getenv("HOME") + "/.teatime/"

const TEATIME_NET_SYN = "teatime_syn"
const TEATIME_NET_ACK = "teatime_ack"
const TEATIME_NET_SYNACK = "teatime_synack"

const TEATIME_ALIVE_PING = "tt_alive?"
const TEATIME_GUCCI_PONG = "tt_we_gucci"

// Utils

func ReadFile(path string) ([]string, error) {
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

func ByteArrayStringEquals(a []byte, b string) bool {
	n := bytes.IndexByte(a, 0)
	return n <= len(b) && string(a[:n]) == b
}

// File object struct definition.  Used for diffs.
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

func (f *File) ToString() string {
	return strings.Join(f.lineSlice, "\n")
}

func (f *File) FromString(s string) {
	f.lineSlice = strings.Split(s, "\n")
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
		_, err = writer.WriteString(fileObj.GetLine(i) + "\n")
		if err != nil {
			return err
		}
	}
	writer.Flush()
	return err
}
