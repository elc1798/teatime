package p2p

import (
	"net"
)

func (this *TTNetSession) StartCrumpetWatcher(socketPath string) error {
	unixAddr, err := net.ResolveUnixAddr("unix", socketPath)
	if err != nil {
		return err
	}

	conn, err := net.DialUnix("unix", nil, unixAddr)
	if err != nil {
		return err
	}

	this.CrumpetWatcher = conn
	go this.watchCrumpet()

	return nil
}

func (this *TTNetSession) watchCrumpet() {
	// TODO: Do things based on serialized objects that are passed around
}
