package p2p

import (
    "log"
    "net"
)

type Peer struct {
    IP      string
    Port    int
}

type TTNetSession struct {
    CAConn      *net.TCPConn    // Connection to central authority
    PeerConns   []*net.TCPConn  // List of peer connections
    PeerList    []Peer
}

const CentralAuthorityHost = "tsukiumi.elc1798.tech:9001"
const PeerCacheFile = "~/.teatime/peer_cache"

func NewTTNetSession() (*TTNetSession) {
    newSession := new(TTNetSession)
    newSession.CAConn = nil
    newSession.PeerConns = make([]*net.TCPConn, 0)
    newSession.PeerList = make([]Peer, 0)

    return newSession
}

func makeTCPConn(host string) (*net.TCPConn, error) {
    tcpAddr, err := net.ResolveTCPAddr("tcp", host)
    if err != nil {
        return nil, err
    }

    return net.DialTCP("tcp", nil, tcpAddr)
}

/*
 * Attempts to make a TCP connection to the designated destination. If a TCP
 * connection is made, reports to central authority that we are now part of the
 * peer network.
 *
 * 'host' should be a string in the form "<server>:<port>" where <server> is an
 * IPv4 address or a domain name, and port is a valid network port.
 */
func (this *TTNetSession) TryTeaTimeConn(host string) (error) {
    peerConnection, err := makeTCPConn(host)
    if err != nil {
        return err
    }

    // Connect to central authority
    if this.CAConn == nil {
        this.CAConn, err = makeTCPConn(CentralAuthorityHost)

        // Errors to central authority on non-fatal
        if err != nil {
            log.Println(err)
            this.CAConn = nil
        }
    }

    if this.CAConn != nil {
        _, err = this.CAConn.Write([]byte(host))
        reply := make([]byte, 1024)
        _, err = this.CAConn.Read(reply)

        if string(reply) != "ok" {
            log.Printf("Error from server reply: (err) %v, (resp) %v", err, reply)
        }
    }

    // Add peer to internal tracking
    this.PeerConns = append(this.PeerConns, peerConnection)
    tcpAddr, _ := net.ResolveTCPAddr("tcp", host)
    this.PeerList = append(this.PeerList, Peer{
        IP: string(tcpAddr.IP),
        Port: tcpAddr.Port,
    })

    return nil
}

/*
 * Gets a list of peers from local cache
 */
func GetLocalPeerCache() ([]Peer) {
    return make([]Peer, 0)
}

/*
 * Gets a list of peers from central authority
 */
func GetPeerListFromCentral() ([]Peer, error) {
    return make([]Peer, 0), nil
}

