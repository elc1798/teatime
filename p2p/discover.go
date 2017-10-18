package p2p

import (
    "net"
)

type Peer struct {
    IP      string
    Port    int
    ID      string
}

/*
 * Attempts to make a TCP connection to the designated destination. If a TCP
 * connection is made, reports to central authority that we are now part of the
 * peer network.
 *
 * Returns the TCPConn object, if connection was made. nil if unsuccessful.
 */
func TryTeaTimeConn(network string, laddr, raddr *net.TCPAddr) (*net.TCPConn, error) {
    return nil, nil
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

