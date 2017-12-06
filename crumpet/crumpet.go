package crumpet

import (
	"log"
	"net"
	"os"

	tt "github.com/elc1798/teatime"
	fs "github.com/elc1798/teatime/fs"
	p2p "github.com/elc1798/teatime/p2p"
)

type CrumpetDaemon struct {
	impendingConnections map[string]chan bool
	repoSockets          map[string]*net.UnixConn
	peerConnections      map[string][]string
	netSessions          map[string]*p2p.TTNetSession
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
	this.startCLISocket()
}

func waitForInput(c chan bool, s string) {
	<-c
	log.Printf("Repo '%v' connected to Crumpet!", s)
}

func (this *CrumpetDaemon) initFields() error {
	// Set up socket directory
	os.RemoveAll(tt.TEATIME_SOCKET_DIR)
	if err := os.MkdirAll(tt.TEATIME_SOCKET_DIR, 0755); err != nil {
		return err
	}

	// Set up Teatime directory
	os.MkdirAll(tt.TEATIME_DEFAULT_HOME, 0755)
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
	this.netSessions = make(map[string]*p2p.TTNetSession)

	repoList, err := fs.GetAllRepos()
	if err != nil {
		return err
	}

	for _, repo := range repoList {
		if e1 := this.setUpRepoRoutines(repo); e1 != nil {
			return e1
		}
	}

	return nil
}

func (this *CrumpetDaemon) setUpRepoRoutines(repo *fs.Repo) error {
	socketPath := tt.GetSocketPath(repo.Name)

	listener, err := startUnixSocket(socketPath)
	if err != nil {
		return err
	}

	waitChan := make(chan bool)
	go waitForInput(waitChan, repo.Name)
	go this.waitForNetSession(listener, repo, waitChan)
	this.impendingConnections[repo.Name] = waitChan

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

	listener.SetUnlinkOnClose(false)
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
