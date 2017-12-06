package main

import (
	"log"
	"net"

	tt "github.com/elc1798/teatime"
	crumpet "github.com/elc1798/teatime/crumpet"
	encoder "github.com/elc1798/teatime/encode"
)

func startCrumpetAndHang() {
	daemon := crumpet.NewCrumpetDaemon()
	daemon.Start(true)

	c := make(chan bool)
	<-c
}

func sendCrumpetCommand(args []string) {
	s1 := encoder.IntraTeatimeSerializer{}
	encoded, err := s1.Serialize(args)
	if err != nil {
		log.Printf("Failed to serialize data!")
		return
	}

	unixAddr, _ := net.ResolveUnixAddr("unix", tt.TEATIME_CLI_SOCKET)
	conn, err := net.DialUnix("unix", nil, unixAddr)
	if err != nil {
		log.Printf("Failed to connect to Crumpet")
		return
	}
	defer conn.Close()

	_, err = tt.SendData(conn, encoded)
	if err != nil {
		log.Printf("Failed to send data to Crumpet")
		return
	}

	resp, _, err := tt.ReadData(conn)
	if err != nil {
		log.Printf("Failed to get response from Crumpet")
		return
	}

	log.Print(string(resp))
}
