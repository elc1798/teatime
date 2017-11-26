package crumpet

import (
	"errors"
	"log"
	"net"
	"strconv"
	"time"

	tt "github.com/elc1798/teatime"
	encoder "github.com/elc1798/teatime/encode"
)

func (this *CrumpetDaemon) StartListener(global bool) error {
	if this.Listener != nil {
		return errors.New("Crumpet.Listener: Already started")
	}

	IP := "127.0.0.1:"
	if global {
		IP = "0.0.0.0:"
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", IP+strconv.Itoa(tt.TEATIME_DEFAULT_PORT))
	if err != nil {
		return nil
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}

	this.Listener = listener
	go this.listenerAcceptLoop()
	return nil
}

func (this *CrumpetDaemon) listenerAcceptLoop() {
	defer this.Listener.Close()

	for {
		conn, err := this.Listener.AcceptTCP()
		if err != nil {
			log.Printf("Crumpet.Listener: Error accepting connection: %v\n", err)
			continue
		}

		log.Println("Crumpet.Listener: Accepting connection")
		go this.handleConnection(conn)
	}
}

func (this *CrumpetDaemon) handleConnection(conn *net.TCPConn) {
	// Do Teatime handshake
	desiredRepo, err := waitForConnectionRequest(conn)
	if err != nil {
		conn.Close()
		log.Printf("Crumpet.Listener: HandshakeError->(%v)", err)
		return
	}

	// Check if Repo is connected
	if _, ok := this.repoSockets[desiredRepo]; !ok {
		conn.Close()
		log.Printf("Crumpet.Listener: Repo '%v' not connected to Crumpet!", desiredRepo)
		return
	}

	// Get connection host
	host := conn.RemoteAddr().String()
	tcpAddr, _ := net.ResolveTCPAddr("tcp", host)
	log.Printf("Crumpet.Listener: NewPeerConnection->(%v -> %v)", tcpAddr.IP, desiredRepo)

	// Serialize the IP of the peer and relay to appropriate TTNetSession
	serializer := encoder.InterTeatimeSerializer{}
	connectionInfo := encoder.TeatimeMessage{
		Recipient: desiredRepo,
		Action:    encoder.ACTION_CONNECT,
		Payload:   encoder.ConnectionRequestPayload(tcpAddr.IP),
	}

	// This really should not error...
	encoded, _ := serializer.Serialize(connectionInfo)
	this.repoSockets[desiredRepo].Write(encoded)
}

func waitForConnectionRequest(conn *net.TCPConn) (string, error) {
	conn.SetReadDeadline(time.Now().Add(time.Second * 2))
	data, _, err := tt.ReadData(conn)
	if err != nil {
		return "", err
	}

	serializer := encoder.InterTeatimeSerializer{}
	decoded_obj, err := serializer.Deserialize(data)
	if err != nil {
		return "", err
	}

	decoded, ok := decoded_obj.(encoder.TeatimeMessage)
	if !ok {
		return "", errors.New("Invalid TeatimeMessage")
	}

	if decoded.Action != encoder.ACTION_CONNECT {
		return "", errors.New("Not a connection attempt")
	}

	return decoded.Recipient, nil
}
