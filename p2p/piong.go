/*
 * This file defines functions that are called by peers to communicate changes
 * and differences. The file is called "piong" as a shortening of 'ping pong',
 * as data will be constantly be sent to and from peers.
 */

package p2p

import (
	"log"
	"time"

	tt "github.com/elc1798/teatime"
)

func (this *TTNetSession) sendTTPing(peerID string) error {
	// Send "alive" signal to check if peer is up
	this.NumPingsSent++
	_, err := SendData(this.PeerConns[peerID], []byte(tt.TEATIME_ALIVE_PING))
	return err
}

func (this *TTNetSession) sendTTPong(peerID string) error {
	// Send "we gucci" signal to check if peer is up
	this.NumPongsSent++
	_, err := SendData(this.PeerConns[peerID], []byte(tt.TEATIME_GUCCI_PONG))
	return err
}

func (this *TTNetSession) checkPeer(peerID string) bool {
	this.sendTTPing(peerID)

	this.PeerConns[peerID].SetReadDeadline(time.Now().Add(time.Second * 2))
	if resp, _, err := ReadData(this.PeerConns[peerID]); err != nil || !tt.ByteArrayStringEquals(resp, tt.TEATIME_GUCCI_PONG) {
		log.Printf("Check peer failed: err=%v, data=%v", err, string(resp))

		// Close connection
		this.PeerConns[peerID].Close()

		// Remove peer
		delete(this.PeerConns, peerID)
		delete(this.PeerList, peerID)

		return false
	} else {
		this.NumPongsRcvd++

		return true
	}
}

func (this *TTNetSession) respondToData(peerID string) error {
	if resp, _, err := ReadData(this.PeerConns[peerID]); err != nil {
		return err
	} else {
		if tt.ByteArrayStringEquals(resp, tt.TEATIME_ALIVE_PING) {
			this.NumPingsRcvd++

			// Send we gucci
			this.sendTTPong(peerID)
		} else {
			// TODO: Send received to delegate function
			log.Printf("Received: %v", string(resp))
		}
	}

	return nil
}

func (this *TTNetSession) startConnListener(peerID string) {
	for {
		if err := this.respondToData(peerID); err != nil {
			log.Printf("Connection responder crashed: %v", err)
			break
		}
	}
}

func (this *TTNetSession) startPingService(peerID string, pingInterval time.Duration) {
	ticker := time.NewTicker(pingInterval)

	for {
		log.Printf("Pinging %v at %v", peerID, time.Now())
		if !this.checkPeer(peerID) {
			log.Printf("Peer %v failed. Stopping pings", peerID)
			ticker.Stop()
			break
		}

		<-ticker.C
	}
}

func (this *TTNetSession) startChangeNotifier(peerID string) {
	ticker := time.NewTicker(1 * time.Second)

	for {
		log.Printf("Polling for new files at %v", time.Now())
		changedFiles, err := this.Repo.GetChangedFiles()

		log.Printf("Files changed: %v, err: %v", changedFiles, err)
		if changedFiles != nil && len(changedFiles) > 0 {
			s := ChangedFileListSerializer{}
			encoded, err := s.Serialize(changedFiles)
			if err == nil {
				SendData(this.PeerConns[peerID], encoded)
			}
		}
		<-ticker.C
	}
}
