package p2p

import (
	"fmt"
	"log"
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
	if _, ok := this.PeerList[ipv4]; ok {
		return fmt.Errorf("TTNetSession.Discover: Peer '%v' already exists", ipv4)
	}

	log.Printf("Repo %v sending connection request to %v", this.Repo.Name, repoName)
	peerConnection, err := makeTCPConn(host)
	if err != nil {
		return err
	}

	if e1 := this.sendConnectionRequest(peerConnection, repoName); e1 != nil {
		return e1
	}
	peerConnection.Close()
	log.Printf("Repo %v sent connection request to %v", this.Repo.Name, repoName)

	return nil
}
