package crumpet

import (
	"fmt"
	"log"
	"net"

	tt "github.com/elc1798/teatime"
	encoder "github.com/elc1798/teatime/encode"
)

// When Crumpet receives packed data, deserialize it. Data should contain
// desired repo, and the data to be passed to that repo. (Use
// encode.Serializer!) Crumpet should delegate to the corresponding repoSocket.

func (this *CrumpetDaemon) Delegate(conn *net.TCPConn) {
	defer conn.Close()

	// Assumes that connection was made already
	data, _, err := tt.ReadData(conn)
	if err != nil {
		log.Printf("Crumpet.Delegator: Error reading data: %v", err)
		log.Printf("Data gotten so far: %v", string(data))
		log.Printf("As array: %v", data)
		return
	}

	serializer := encoder.InterTeatimeSerializer{}
	decoded_obj, err := serializer.Deserialize(data)
	if err != nil {
		log.Printf("Crumpet.Delegator: Error deserializing data: %v", err)
		return
	}

	decoded, ok := decoded_obj.(encoder.TeatimeMessage)
	if !ok {
		return
	}

	switch decoded.Action {
	case encoder.ACTION_CONNECT:
		if err := this.handleActionConnReq(conn, decoded); err != nil {
			log.Printf("Crumpet.Delegator.ActionConnReq error: %v", err)
		}
	case encoder.ACTION_PING:
		if err := this.handleActionPing(conn, decoded); err != nil {
			log.Printf("Crumpet.Delegator.ActionPing error: %v", err)
		}
	case encoder.ACTION_DELTAS:
		if err := this.handleActionFileDeltas(conn, decoded); err != nil {
			log.Printf("Crumpet.Delegator.ActionFileDeltas error: %v", err)
		}
	}
}

func repoNotConnectedError(repoName string) error {
	return fmt.Errorf("Crumpet.Delegator: Repo '%v' not connected to Crumpet!", repoName)
}

func getIPFromRemoteConn(conn *net.TCPConn) string {
	host := conn.RemoteAddr().String()
	tcpAddr, _ := net.ResolveTCPAddr("tcp", host)
	return fmt.Sprintf("%v", tcpAddr.IP)
}

func (this *CrumpetDaemon) handleActionConnReq(conn *net.TCPConn, msg encoder.TeatimeMessage) error {
	desiredRepo := msg.Recipient

	// Check if Repo is connected
	if _, ok := this.repoSockets[desiredRepo]; !ok {
		return repoNotConnectedError(desiredRepo)
	}

	connectInfo, ok := msg.Payload.(encoder.ConnectionRequestPayload)
	if !ok {
		return fmt.Errorf("Crumpet.Delegator: Invalid conn_req payload")
	}

	// Get connection host
	originIP := getIPFromRemoteConn(conn)
	log.Printf("Crumpet.Delegator: NewPeerConnection->(%v -> %v)", originIP, desiredRepo)

	// Serialize the IP of the peer and relay to appropriate TTNetSession
	serializer := encoder.InterTeatimeSerializer{}
	connectionInfo := encoder.TeatimeMessage{
		Recipient: desiredRepo,
		Action:    encoder.ACTION_CONNECT,
		Payload: encoder.ConnectionRequestPayload{
			OriginIP:       originIP,
			RepoRemoteName: connectInfo.RepoRemoteName,
		},
	}

	// This really should not error...
	encoded, _ := serializer.Serialize(connectionInfo)
	if _, err := tt.SendData(this.repoSockets[desiredRepo], encoded); err != nil {
		return err
	}

	// Add to peerConnections
	currConns, ok := this.peerConnections[originIP]
	if ok {
		this.peerConnections[originIP] = append(currConns, desiredRepo)
	} else {
		this.peerConnections[originIP] = []string{desiredRepo}
	}

	return nil
}

func (this *CrumpetDaemon) handleActionPing(conn *net.TCPConn, msg encoder.TeatimeMessage) error {
	pingInfo, ok := msg.Payload.(encoder.PingPayload)
	if !ok {
		return fmt.Errorf("Crumpet.Delegator: Invalid ping payload")
	}

	// Check if Repo is connected
	if _, ok := this.repoSockets[msg.Recipient]; !ok {
		return repoNotConnectedError(msg.Recipient)
	}

	// Repack the message and send to appropriate TTNetSession
	pingInfo.OriginIP = getIPFromRemoteConn(conn)
	msg.Payload = pingInfo

	serializer := encoder.InterTeatimeSerializer{}
	encoded, _ := serializer.Serialize(msg)

	log.Printf("Crumpet.HandleActionPing: encoded=%v", string(encoded))
	_, err := tt.SendData(this.repoSockets[msg.Recipient], encoded)
	return err
}

func (this *CrumpetDaemon) handleActionFileDeltas(conn *net.TCPConn, msg encoder.TeatimeMessage) error {
	deltasInfo, ok := msg.Payload.(encoder.FileDeltasPayload)
	if !ok {
		return fmt.Errorf("Crumpet.Delegator: Invalid file deltas payload")
	}

	// Check if Repo is connected
	if _, ok := this.repoSockets[msg.Recipient]; !ok {
		return repoNotConnectedError(msg.Recipient)
	}

	deltasInfo.OriginIP = getIPFromRemoteConn(conn)
	msg.Payload = deltasInfo

	serializer := encoder.InterTeatimeSerializer{}
	encoded, _ := serializer.Serialize(msg)

	log.Printf("Crumpet.HandleActionFileDeltas: encoded=%v", string(encoded))
	_, err := tt.SendData(this.repoSockets[msg.Recipient], encoded)
	return err
}
