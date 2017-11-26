package p2p

/*
 * These functions are only here for TESTING purposes. They should not be called
 * OUTSIDE of tests.
 */

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	tt "github.com/elc1798/teatime"
	fs "github.com/elc1798/teatime/fs"
)

type TestSession struct {
	*TTNetSession

	Listener *net.TCPListener
}

/*
 * Performs a 3-way handshake (syn, ack, syn-ack) from the server's perspective
 */
func doTeatimeServerHandshake(conn *net.TCPConn) error {
	// Confirm syn
	conn.SetReadDeadline(time.Now().Add(time.Second * 2)) // Stop read after 2 seconds
	if remote_syn, _, err := tt.ReadData(conn); err != nil || !tt.ByteArrayStringEquals(remote_syn, tt.TEATIME_NET_SYN) {
		return fmt.Errorf("Invalid teatime_syn [%v]", err)
	}

	// Send ack
	if _, err := tt.SendData(conn, []byte(tt.TEATIME_NET_ACK)); err != nil {
		return fmt.Errorf("Error sending ack [%v]", err)
	}

	// Wait for syn-ack
	conn.SetReadDeadline(time.Now().Add(time.Second * 2)) // Stop read after 2 seconds
	if remote_synack, _, err := tt.ReadData(conn); err != nil || !tt.ByteArrayStringEquals(remote_synack, tt.TEATIME_NET_SYNACK) {
		return fmt.Errorf("Invalid teatime_synack [%v]", err)
	}

	return nil
}

/*
 * Handles an incoming TCP connection. Performs handshake and metadata
 * bookkeeping
 */
func handleConnection(sess *TestSession, conn *net.TCPConn) {
	// Handle syn-ack handshake
	if err := doTeatimeServerHandshake(conn); err != nil {
		conn.Close()
		fmt.Sprintf("Handshake error: < %v >", err)
		return
	}

	// Get connection host
	host := conn.RemoteAddr().String()
	tcpAddr, _ := net.ResolveTCPAddr("tcp", host)
	key := fmt.Sprintf("%v:%v", tcpAddr.IP, tcpAddr.Port)

	// Add to peer conns
	sess.PeerConns[key] = conn

	// Add to peer list
	sess.PeerList[key] = Peer{
		IP:   fmt.Sprintf("%v", tcpAddr.IP),
		Port: tcpAddr.Port,
	}

	newPeer := sess.PeerList[key]
	log.Printf("Listener: [New Connected Peer] %v:%v", newPeer.IP, newPeer.Port)

	// Start listener
	go sess.startConnListener(key)
	sess.GenerateLocalPeerCache()
}

/*
 * Runs continuous loop for listener to accept connections
 */
func (this *TestSession) listenerAcceptLoop() {
	defer this.Listener.Close()

	for {
		conn, err := this.Listener.AcceptTCP()
		if err != nil {
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}

		log.Println("Listener: Accepting connection")
		go handleConnection(this, conn)
	}
}

/*
 * Starts the listener. Setting `global` to `true` will allow connections to be
 * accepted from non-local addresses.
 */
func (this *TestSession) StartListener(port int, global bool) error {
	if this.Listener != nil {
		return errors.New("Listener: Already started")
	}

	IP := "127.0.0.1:"
	if global {
		IP = "0.0.0.0:"
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", IP+strconv.Itoa(port))
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

func NewTestSession(repo *fs.Repo) *TestSession {
	sess := NewTTNetSession(repo)
	return &TestSession{sess, nil}
}
