package teatime

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"path"
	"strings"
)

// Project constants

const TEATIME_TRACKED_DIR = ".tracked/"
const TEATIME_BACKUP_DIR = ".backup/"
const TEATIME_PEER_CACHE = "/peer_cache"
const TEATIME_SOCKET_DIR = "/tmp/teatime/"
const TEATIME_DIR_ROOT_STORE = "/dir_root"

var TEATIME_DEFAULT_HOME = path.Join(os.Getenv("HOME"), ".teatime/")

const TEATIME_DEFAULT_PORT = 12345

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

func GetSocketPath(repoName string) string {
	return path.Join(TEATIME_SOCKET_DIR, fmt.Sprintf("%v.sock", repoName))
}

func ResetTeatime() error {
	os.RemoveAll(TEATIME_DEFAULT_HOME)
	os.RemoveAll(TEATIME_SOCKET_DIR)
	return os.Mkdir(TEATIME_DEFAULT_HOME, 0755)
}

// TCP connection helpers

/*
 * Send specified data to the specified connection
 *
 * Returns the number of bytes sent, and an error if unsuccessful
 */
func SendData(conn *net.TCPConn, bytes []byte) (int, error) {
	return conn.Write(append(bytes, byte(0)))
}

/*
 * Reads data from the specified connection
 *
 * Returns a byte array containing the read data, number of bytes read and an
 * error if unsuccessful
 */
func ReadData(conn *net.TCPConn) ([]byte, int, error) {
	full_reply := make([]byte, 0)
	reply := make([]byte, 1)
	reply[0] = byte(1) // Doesn't matter as long as it's not 0
	for {
		num_bytes, err := conn.Read(reply)
		if err != nil {
			return full_reply, len(full_reply), err
		}

		if num_bytes != 1 {
			continue
		}

		if reply[0] == byte(0) {
			break
		}

		full_reply = append(full_reply, reply[0])
	}
	return full_reply, len(full_reply), nil
}
