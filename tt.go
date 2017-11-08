package teatime

import (
	"bufio"
	"bytes"
	"os"
)

// Project constants

const TEATIME_TRACKED_DIR = ".tracked/"
const TEATIME_BACKUP_DIR = ".backup/"

var TEATIME_DEFAULT_HOME = os.Getenv("HOME") + "/.teatime"
var TEATIME_PEER_CACHE = TEATIME_DEFAULT_HOME + "/peer_cache"

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
