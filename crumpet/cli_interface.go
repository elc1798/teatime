package crumpet

import (
	"fmt"
	"log"
	"net"

	tt "github.com/elc1798/teatime"
	encoder "github.com/elc1798/teatime/encode"
	fs "github.com/elc1798/teatime/fs"
	p2p "github.com/elc1798/teatime/p2p"
)

func logAndSendMsg(message string, conn *net.UnixConn, close bool) {
	log.Printf(message)
	tt.SendData(conn, []byte(message))

	if close {
		conn.Close()
	}
}

func (this *CrumpetDaemon) startCLISocket() error {
	unixAddr, err := net.ResolveUnixAddr("unix", tt.TEATIME_CLI_SOCKET)
	if err != nil {
		return err
	}

	listener, err := net.ListenUnix("unix", unixAddr)
	if err != nil {
		return err
	}

	listener.SetUnlinkOnClose(true)
	go this.waitForCLIConnection(listener)

	return nil
}

func (this *CrumpetDaemon) waitForCLIConnection(listener *net.UnixListener) {
	defer listener.Close()

	for {
		conn, err := listener.AcceptUnix()
		if err != nil {
			continue
		}

		// Read from conn. Should contain a command
		data, _, err := tt.ReadData(conn)
		if err != nil {
			logAndSendMsg(fmt.Sprintf("Crumpet.CLI: Error reading data: %v", err), conn, true)
			continue
		}

		serializer := encoder.IntraTeatimeSerializer{}
		decoded_obj, err := serializer.Deserialize(data)
		if err != nil {
			logAndSendMsg(fmt.Sprintf("Crumpet.CLI: Error decoding data: %v", err), conn, true)
			continue
		}

		decoded, ok := decoded_obj.([]string)
		if !ok {
			logAndSendMsg(fmt.Sprintf("Crumpet.CLI: Error converting JSON to object"), conn, true)
			continue
		}

		this.handleCLICommand(conn, decoded)
	}
}

func (this *CrumpetDaemon) handleCLICommand(conn *net.UnixConn, decoded []string) {
	switch decoded[0] {
	case encoder.COMMAND_INIT_REPO:
		if len(decoded) != 3 {
			logAndSendMsg(fmt.Sprintf("Crumpet.CLI: Invalid INIT_REPO command: %v", decoded), conn, true)
			return
		}

		new_repo, err := fs.InitRepo(decoded[1], decoded[2])
		if err != nil {
			logAndSendMsg(fmt.Sprintf("Crumpet.CLI: Error creating repo: %v", err), conn, true)
			return
		}

		if e1 := this.setUpRepoRoutines(new_repo); e1 != nil {
			logAndSendMsg(fmt.Sprintf("Error starting repo network listeners: %v", err), conn, true)
			return
		}

		// Start the TeaTimeSession
		this.netSessions[new_repo.Name] = p2p.NewTTNetSession(new_repo)
	case encoder.COMMAND_LINK_PEER:
		if len(decoded) != 4 {
			logAndSendMsg(fmt.Sprintf("Crumpet.CLI: Invalid LINK_PEER command: %v", decoded), conn, true)
			return
		}

		if _, ok := this.netSessions[decoded[1]]; !ok {
			logAndSendMsg(fmt.Sprintf("Crumpet.CLI: Repo %v does not exist", decoded[1]), conn, true)
			return
		}

		if e1 := this.netSessions[decoded[1]].TryTeaTimeConn(decoded[2], decoded[3]); e1 != nil {
			logAndSendMsg(fmt.Sprintf("Crumpet.CLI: Could not connect to peer: %v", e1), conn, true)
			return
		}
	case encoder.COMMAND_ADD_FILE:
	default:
		logAndSendMsg(fmt.Sprintf("Crumpet.CLI: Invalid command: %v", decoded[0]), conn, true)
	}

	logAndSendMsg("ok", conn, true)
}
