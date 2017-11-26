package p2p

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strings"
	"time"

	tt "github.com/elc1798/teatime"
	encoder "github.com/elc1798/teatime/encode"
	fs "github.com/elc1798/teatime/fs"
)

type Peer struct {
	IP   string `json:"IP"`
	Port int    `json:"Port"`
}

type TTNetSession struct {
	CAConn         *net.TCPConn            // Connection to central authority
	PeerConns      map[string]*net.TCPConn // List of peer connections
	PeerList       map[string]Peer
	Repo           *fs.Repo
	CrumpetWatcher *net.UnixConn // Unix Domain Socket for Crumpet communication

	// Counters
	NumPingsSent int
	NumPingsRcvd int
	NumPongsSent int
	NumPongsRcvd int
}

var peerListSerializer = encoder.NewDefaultSerializer(make(map[string]Peer))

/*
 * Creates and initializes a new Teatime Network Session
 */
func NewTTNetSession(repo *fs.Repo) *TTNetSession {
	newSession := new(TTNetSession)
	newSession.Repo = repo
	newSession.CAConn = nil
	newSession.PeerConns = make(map[string]*net.TCPConn)
	p_list, err := newSession.GetLocalPeerCache()
	if err != nil {
		log.Printf("NewTTNetSession: Error in GetLocalPeerCache->[ %v ]", err)
		p_list = make(map[string]Peer)
	}
	newSession.PeerList = p_list

	for _, peer := range newSession.PeerList {
		_ = newSession.TryTeaTimeConn(fmt.Sprintf("%s:%d", peer.IP, peer.Port), time.Millisecond*250)
	}

	newSession.NumPingsSent = 0
	newSession.NumPingsRcvd = 0
	newSession.NumPongsSent = 0
	newSession.NumPingsRcvd = 0

	return newSession
}

/*
 * Gets a list of peers from local cache
 */
func (this *TTNetSession) GetLocalPeerCache() (map[string]Peer, error) {
	peer_data, err := tt.ReadFile(this.Repo.GetPeerCacheFile())
	if err != nil {
		return nil, err
	}

	generic_obj, err := peerListSerializer.Deserialize([]byte(strings.Join(peer_data, "\n")))
	if err != nil {
		return nil, err
	}

	// Recast
	peer_list := make(map[string]Peer)
	for k, v := range generic_obj.(map[string]interface{}) {
		p := v.(map[string]interface{})
		peer_list[k] = Peer{
			IP:   p["IP"].(string),
			Port: int(p["Port"].(float64)),
		}
	}

	return peer_list, nil
}

func (this *TTNetSession) GenerateLocalPeerCache() error {
	bytes, err := peerListSerializer.Serialize(this.PeerList)
	if err != nil {
		return err
	}

	// Write to file
	return ioutil.WriteFile(
		this.Repo.GetPeerCacheFile(),
		bytes,
		0644,
	)
}