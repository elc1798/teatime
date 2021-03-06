package p2p

import (
	"errors"
	"log"
	"net"
	"time"

	tt "github.com/elc1798/teatime"
	encoder "github.com/elc1798/teatime/encode"
)

var PING_INTERVAL = time.Millisecond * 5000

func (this *TTNetSession) startCrumpetWatcher() error {
	socketPath := tt.GetSocketPath(this.Repo.Name)
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
	for {
		crumpetData, _, err := tt.ReadData(this.CrumpetWatcher)
		if err != nil {
			log.Printf("Error reading from Crumpet!")
			continue
		}

		serializer := encoder.InterTeatimeSerializer{}
		decoded_obj, err := serializer.Deserialize(crumpetData)
		if err != nil {
			log.Printf("Error decoding Crumpet data: %v", err)
			continue
		}

		decoded, ok := decoded_obj.(encoder.TeatimeMessage)
		if !ok {
			continue
		}

		switch decoded.Action {
		case encoder.ACTION_CONNECT:
			if e2 := this.handleActionConnect(decoded.Payload); e2 != nil {
				log.Printf("HandleActionConnectError: %v", e2)
			}
		case encoder.ACTION_PING:
			if e2 := this.handleActionPing(decoded.Payload); e2 != nil {
				log.Printf("HandleActionPingError: %v", e2)
			}
		case encoder.ACTION_DELTAS:
			if e2 := this.handleActionFileDeltas(decoded.Payload); e2 != nil {
				log.Printf("HandleFileListError: %v", e2)
			}
		}
	}
}

func (this *TTNetSession) handleActionConnect(v interface{}) error {
	connectInfo, ok := v.(encoder.ConnectionRequestPayload)
	if !ok {
		return errors.New("Invalid ConnectionRequestPayload")
	}

	log.Printf("handleActionConnect called: %v", connectInfo)
	peerIP := connectInfo.OriginIP
	if _, ok := this.PeerList[peerIP]; ok {
		// We don't need to send a new connection. This is just the residual
		// connection request from us sending them a connection request.
		log.Printf("TTNetSession.Receiver: Peer '%v' already exists", peerIP)
		return nil
	}
	log.Printf("Repo %v received connection request from %v", this.Repo.Name, connectInfo.RepoRemoteName)

	if err := this.TryTeaTimeConn(peerIP, connectInfo.RepoRemoteName); err != nil {
		// Remove from peer list
		return err
	}

	this.PeerList[peerIP] = Peer{
		IP:             peerIP,
		Port:           tt.TEATIME_DEFAULT_PORT,
		RepoRemoteName: connectInfo.RepoRemoteName,
	}

	log.Printf("Repo %v accepted request from %v", this.Repo.Name, connectInfo.RepoRemoteName)
	// Start the Ping Service after we've connected to them and handled
	// their connection
	go this.startPingService(peerIP, PING_INTERVAL, connectInfo.RepoRemoteName)

	// Start file diff service
	go this.startFileTrackerService(peerIP, connectInfo.RepoRemoteName)

	return nil
}

func (this *TTNetSession) handleActionPing(v interface{}) error {
	pingInfo, ok := v.(encoder.PingPayload)
	if !ok {
		return errors.New("Invalid PingPayload")
	}

	log.Printf("handleActionPing called: %v", pingInfo)
	if _, hasPeer := this.PeerList[pingInfo.OriginIP]; !hasPeer {
		return errors.New("Ping originated from unregistered peer!")
	}

	if pingInfo.IsPong {
		this.NumPongsRcvd++

		return nil
	} else {
		this.NumPingsRcvd++

		return this.sendTTPong(pingInfo.OriginIP)
	}
}

func (this *TTNetSession) handleActionFileDeltas(v interface{}) error {
	deltasInfo, ok := v.(encoder.FileDeltasPayload)
	if !ok {
		return errors.New("Invalid FileDeltasPayload")
	}

	log.Printf("handleActionFileDeltas called: %v", deltasInfo)
	if _, hasPeer := this.PeerList[deltasInfo.OriginIP]; !hasPeer {
		return errors.New("Ping originated from unregistered peer!")
	}

	if deltasInfo.IsAck {
		// Diffs that remote did not apply were removed from map
		for fileName, _ := range deltasInfo.Deltas {
			this.Repo.WriteBackupFile(fileName)
		}

		// Resume polling
		this.resumeFilePollChans[deltasInfo.OriginIP] <- true
		return nil
	}

	// Get our own changed files
	ourChangedFiles, err := this.Repo.GetChangedFiles()
	if err != nil {
		return err
	}
	ourDiffStrings, _ := this.Repo.GetDiffStrings(ourChangedFiles)
	fileDiffMap := make(map[string]string)
	for i, v := range ourChangedFiles {
		fileDiffMap[v] = ourDiffStrings[i]
	}

	appliedFiles := make(map[string]string)

	// We need to check if any of the files they changed correspond with files
	// we changed, then send the appropriate diffs over. Patch the ones that we
	// have not yet changed.
	for fileName, diffString := range deltasInfo.Deltas {
		ourDiff, ok := fileDiffMap[fileName]
		if ok {
			err = this.Repo.PatchFileMergeConflict(fileName, []string{diffString, ourDiff})
			if err != nil {
				log.Printf("%v: Error fixing merge conflict: %v", this.Repo.Name, err)
				continue
			}
		} else {
			err = this.Repo.PatchFile(fileName, diffString)
			if err != nil {
				log.Printf("%v: Error patching: %v", this.Repo.Name, err)
				continue
			}
		}

		this.Repo.WriteBackupFile(fileName)
		appliedFiles[fileName] = "applied"
	}

	// Respond with deltas ack
	this.sendDeltasAck(deltasInfo.OriginIP, appliedFiles)

	return nil
}
