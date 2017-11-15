package p2p

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	tt "github.com/elc1798/teatime"
	fs "github.com/elc1798/teatime/fs"
)

type Peer struct {
	IP   string
	Port int
}

type TTNetSession struct {
	CAConn    *net.TCPConn            // Connection to central authority
	PeerConns map[string]*net.TCPConn // List of peer connections
	PeerList  map[string]Peer
	Listener  *net.TCPListener
	Repo      *fs.Repo

	// Counters
	NumPingsSent int
	NumPingsRcvd int
	NumPongsSent int
	NumPongsRcvd int
}

const CentralAuthorityHost = "tsukiumi.elc1798.tech:9001"

/*
 * Creates and initializes a new Teatime Network Session
 */
func NewTTNetSession(repo *fs.Repo) *TTNetSession {
	newSession := new(TTNetSession)
	newSession.CAConn = nil
	newSession.PeerConns = make(map[string]*net.TCPConn)
	newSession.PeerList = make(map[string]Peer)
	newSession.Listener = nil
	newSession.Repo = repo

	newSession.NumPingsSent = 0
	newSession.NumPingsRcvd = 0
	newSession.NumPongsSent = 0
	newSession.NumPingsRcvd = 0

	return newSession
}

/*
 * Generates a TCP connection
 */
func makeTCPConn(host string) (*net.TCPConn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		return nil, err
	}

	return net.DialTCP("tcp", nil, tcpAddr)
}

/*
 * Equivalent of `doTeatimeServerHandshake`, but from client's perspective
 */
func doTeatimeClientHandshake(conn *net.TCPConn) error {
	if _, err := SendData(conn, []byte(tt.TEATIME_NET_SYN)); err != nil {
		return fmt.Errorf("Error sending syn [%v]", err)
	}

	conn.SetReadDeadline(time.Now().Add(time.Second * 2))
	if serv_ack, _, err := ReadData(conn); err != nil || !tt.ByteArrayStringEquals(serv_ack, tt.TEATIME_NET_ACK) {
		return fmt.Errorf("Invalid teatime_ack [%v]", err)
	}

	if _, err := SendData(conn, []byte(tt.TEATIME_NET_SYNACK)); err != nil {
		return fmt.Errorf("Error sending synack [%v]", err)
	}

	return nil
}

/*
 * Attempts to make a TCP connection to the designated destination. If a TCP
 * connection is made, reports to central authority that we are now part of the
 * peer network.
 *
 * 'host' should be a string in the form "<server>:<port>" where <server> is an
 * IPv4 address or a domain name, and port is a valid network port.
 */
func (this *TTNetSession) TryTeaTimeConn(host string, pingInterval time.Duration) error {
	if this.PeerConns[host] != nil {
		return fmt.Errorf("Peer '%v' already exists", host)
	}

	peerConnection, err := makeTCPConn(host)
	if err != nil {
		return err
	}

	// Verify
	if e1 := doTeatimeClientHandshake(peerConnection); e1 != nil {
		peerConnection.Close()
		return e1
	}

	// Add peer to internal tracking
	tcpAddr, _ := net.ResolveTCPAddr("tcp", host)
	key := fmt.Sprintf("%v:%v", tcpAddr.IP, tcpAddr.Port)

	this.PeerConns[key] = peerConnection
	this.PeerList[key] = Peer{
		IP:   fmt.Sprintf("%v", tcpAddr.IP),
		Port: tcpAddr.Port,
	}

	newPeer := this.PeerList[key]
	log.Printf("TryTeaTimeConn: %v:%v", newPeer.IP, newPeer.Port)

	// Connect to central authority {{{
	if this.CAConn == nil {
		this.CAConn, err = makeTCPConn(CentralAuthorityHost)

		// Errors to central authority are non-fatal
		if err != nil {
			log.Println(err)
			this.CAConn = nil
		}
	}

	if this.CAConn != nil {
		// TODO: Auth, then grab peer list from CA
	}
	// }}}

	// Start ping service
	go this.startPingService(key, pingInterval)

	// Start file difference service
	// go this.startChangeNotifier(key)

	return nil
}

/*
 * Gets a list of peers from local cache
 */
func (this *TTNetSession) GetLocalPeerCache() (map[string]Peer, error) {
	peer_data, err := tt.ReadFile(this.Repo.GetPeerCacheFile())
	if err != nil {
		return nil, err
	}

	peer_list := make(map[string]Peer)
	for _, peer_str := range peer_data {
		trimmed := strings.TrimSpace(peer_str)
		tokens := strings.Split(trimmed, ":")

		if len(tokens) != 2 {
			return nil, fmt.Errorf("Invalid peer: '%v'", trimmed)
		}

		port, err := strconv.Atoi(tokens[1])
		if err != nil {
			return nil, fmt.Errorf("Invalid peer: '%v'", trimmed)
		}
		peer_list[trimmed] = Peer{
			IP:   tokens[0],
			Port: port,
		}
	}

	return peer_list, nil
}

func (this *TTNetSession) GenerateLocalPeerCache(peers map[string]Peer) error {
	// Generate string list from peers
	string_list := make([]string, 0)
	for _, peer := range peers {
		string_list = append(string_list, fmt.Sprintf("%s:%d", peer.IP, peer.Port))
	}

	// Write to file
	return ioutil.WriteFile(
		this.Repo.GetPeerCacheFile(),
		[]byte(strings.Join(string_list, "\n")),
		0644,
	)
}

/*
 * Gets a list of peers from central authority
 */
func GetPeerListFromCentral() ([]Peer, error) {
	return make([]Peer, 0), nil
}
