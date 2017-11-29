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
	msg, err := waitForConnectionRequest(conn)
	if err != nil {
		conn.Close()
		log.Printf("Crumpet.Listener: HandshakeError->(%v)", err)
		return
	}

	if err2 := this.handleActionConnReq(conn, msg); err2 != nil {
		conn.Close()
		log.Printf("Crumpet.Listener: ConnectionRequestError->(%v)", err2)
		return
	}

	// Start Delegator
	go this.StartDelegator(conn)
}

func waitForConnectionRequest(conn *net.TCPConn) (encoder.TeatimeMessage, error) {
	conn.SetReadDeadline(time.Now().Add(time.Second * 2))
	data, _, err := tt.ReadData(conn)
	if err != nil {
		return encoder.TeatimeMessage{}, err
	}

	serializer := encoder.InterTeatimeSerializer{}
	decoded_obj, err := serializer.Deserialize(data)
	// log.Printf("%v, %v", string(data), decoded_obj)
	if err != nil {
		return encoder.TeatimeMessage{}, err
	}

	decoded, ok := decoded_obj.(encoder.TeatimeMessage)
	if !ok {
		return encoder.TeatimeMessage{}, errors.New("Invalid TeatimeMessage")
	}

	if decoded.Action != encoder.ACTION_CONNECT {
		return encoder.TeatimeMessage{}, errors.New("Not a connection attempt")
	}

	return decoded, nil
}
