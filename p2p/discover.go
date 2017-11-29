package p2p

import (
	"fmt"
	"net"

	tt "github.com/elc1798/teatime"
	encoder "github.com/elc1798/teatime/encode"
)

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

func (this *TTNetSession) sendConnectionRequest(conn *net.TCPConn, repoName string) error {
	msg := encoder.TeatimeMessage{
		Recipient: repoName,
		Action:    encoder.ACTION_CONNECT,
		Payload: encoder.ConnectionRequestPayload{
			OriginIP:       "",
			RepoRemoteName: this.Repo.Name,
		},
	}

	serializer := encoder.InterTeatimeSerializer{}
	encoded, err := serializer.Serialize(msg)
	if err != nil {
		return err
	}

	_, err = tt.SendData(conn, encoded)
	return err
}

func (this *TTNetSession) TryTeaTimeConn(IP string, repoName string) error {
	host := fmt.Sprintf("%s:%d", IP, tt.TEATIME_DEFAULT_PORT)
	tcpAddr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		return err
	}

	ipv4 := fmt.Sprintf("%v", tcpAddr.IP)
	if this.PeerConns[ipv4] != nil {
		return fmt.Errorf("TTNetSession.Discover: Peer '%v' already exists", ipv4)
	}

	peerConnection, err := makeTCPConn(host)
	if err != nil {
		return err
	}

	if e1 := this.sendConnectionRequest(peerConnection, repoName); e1 != nil {
		peerConnection.Close()
		return e1
	}

	this.PeerConns[ipv4] = peerConnection

	return nil
}
