package teatime

import (
	"os"
)

const TEATIME_TRACKED_DIR = ".tracked/"
const TEATIME_BACKUP_DIR = ".backup/"

var TEATIME_DEFAULT_HOME = os.Getenv("HOME") + "/.teatime"
var TEATIME_PEER_CACHE = TEATIME_DEFAULT_HOME + "/peer_cache"
