package crumpet

import (
	"errors"
	"log"
	"net"
	"strconv"

	tt "github.com/elc1798/teatime"
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
		go this.Delegate(conn)
	}
}
