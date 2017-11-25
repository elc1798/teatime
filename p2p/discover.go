package p2p

import (
	"fmt"
	"log"
	"net"
	"time"

	tt "github.com/elc1798/teatime"
)

const CentralAuthorityHost = "tsukiumi.elc1798.tech:9001"

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
	if _, err := tt.SendData(conn, []byte(tt.TEATIME_NET_SYN)); err != nil {
		return fmt.Errorf("Error sending syn [%v]", err)
	}

	conn.SetReadDeadline(time.Now().Add(time.Second * 2))
	if serv_ack, _, err := tt.ReadData(conn); err != nil || !tt.ByteArrayStringEquals(serv_ack, tt.TEATIME_NET_ACK) {
		return fmt.Errorf("Invalid teatime_ack [%v]", err)
	}

	if _, err := tt.SendData(conn, []byte(tt.TEATIME_NET_SYNACK)); err != nil {
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
 * Gets a list of peers from central authority
 */
func GetPeerListFromCentral() ([]Peer, error) {
	return make([]Peer, 0), nil
}
