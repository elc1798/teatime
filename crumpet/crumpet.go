package crumpet

import (
	"fmt"
	"log"
	"net"
	"os"
	"path"

	tt "github.com/elc1798/teatime"
	fs "github.com/elc1798/teatime/fs"
)

type CrumpetDaemon struct {
	impendingConnections map[string]chan bool
	repoSockets          map[string]*net.UnixConn
	Listener             *net.TCPListener
}

func (this *CrumpetDaemon) Start() {
}

func (this *CrumpetDaemon) initFields() error {
	// Set up socket directory
	if err := os.MkdirAll(tt.TEATIME_SOCKET_DIR, 0755); err != nil {
		return err
	}

	if err := this.setUpRepoSocketMap(); err != nil {
		log.Fatalf("Crumpet: Error setting up sockets: %v", err)
	}
	this.Listener = nil

	return nil
}

func (this *CrumpetDaemon) setUpRepoSocketMap() error {
	this.impendingConnections = make(map[string]chan bool)
	this.repoSockets = make(map[string]*net.UnixConn)

	repoList, err := fs.GetAllRepos()
	if err != nil {
		return err
	}

	for _, repo := range repoList {
		socketPath := path.Join(tt.TEATIME_SOCKET_DIR, fmt.Sprintf("%v.sock", tt.TEATIME_SOCKET_DIR, repo.Name))

		listener, err := startUnixSocket(socketPath)
		if err != nil {
			return err
		}

		waitChan := make(chan bool)
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
