package teatime

import (
	"bufio"
	"os"
)

// Project constants

const TEATIME_TRACKED_DIR = ".tracked/"
const TEATIME_BACKUP_DIR = ".backup/"

var TEATIME_DEFAULT_HOME = os.Getenv("HOME") + "/.teatime"
var TEATIME_PEER_CACHE = TEATIME_DEFAULT_HOME + "/peer_cache"

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
