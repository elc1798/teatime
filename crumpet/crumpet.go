package crumpet

import (
	"log"
	"net"
	"os"

	tt "github.com/elc1798/teatime"
	fs "github.com/elc1798/teatime/fs"
)

type CrumpetDaemon struct {
	impendingConnections map[string]chan bool
	repoSockets          map[string]*net.UnixConn
	peerConnections      map[string][]string
	Listener             *net.TCPListener
}

func NewCrumpetDaemon() *CrumpetDaemon {
	return new(CrumpetDaemon)
}

func (this *CrumpetDaemon) Start(global bool) {
	if err := this.initFields(); err != nil {
		log.Fatalf("Failure initializing Crumpet: %v", err)
	}

	this.StartListener(global)
}

func waitForInput(c chan bool, s string) {
	<-c
	log.Printf("Repo '%v' connected to Crumpet!", s)
}

func (this *CrumpetDaemon) initFields() error {
	// Set up socket directory
	if err := os.MkdirAll(tt.TEATIME_SOCKET_DIR, 0755); err != nil {
		return err
	}

	if err := this.setUpRepoSocketMap(); err != nil {
		return err
	}
	this.Listener = nil

	return nil
}

func (this *CrumpetDaemon) setUpRepoSocketMap() error {
	this.impendingConnections = make(map[string]chan bool)
	this.repoSockets = make(map[string]*net.UnixConn)
	this.peerConnections = make(map[string][]string)

	repoList, err := fs.GetAllRepos()
	if err != nil {
		return err
	}

	for _, repo := range repoList {
		socketPath := tt.GetSocketPath(repo.Name)

		listener, err := startUnixSocket(socketPath)
		if err != nil {
			return err
		}

		waitChan := make(chan bool)
		go waitForInput(waitChan, repo.Name)
		go this.waitForNetSession(listener, repo, waitChan)
		this.impendingConnections[repo.Name] = waitChan
	}

	return nil
}

func startUnixSocket(socketPath string) (*net.UnixListener, error) {
	unixAddr, err := net.ResolveUnixAddr("unix", socketPath)
	if err != nil {
		return nil, err
	}

	listener, err := net.ListenUnix("unix", unixAddr)
	if err != nil {
		return nil, err
	}

	listener.SetUnlinkOnClose(true)
	return listener, nil
}

func (this *CrumpetDaemon) waitForNetSession(l *net.UnixListener, r *fs.Repo, c chan bool) {
	defer l.Close()

	for {
		conn, err := l.AcceptUnix()
		if err != nil {
			continue
		}

		delete(this.impendingConnections, r.Name)
		this.repoSockets[r.Name] = conn
		c <- true

		return
	}
}
