/*
 * This file defines functions that are called by peers to communicate changes
 * and differences. The file is called "piong" as a shortening of 'ping pong',
 * as data will be constantly be sent to and from peers.
 */

package p2p

import (
	"fmt"
	"log"
	"time"

	tt "github.com/elc1798/teatime"
	encoder "github.com/elc1798/teatime/encode"
)

func (this *TTNetSession) sendTTPing(peerID string) error {
	msg := encoder.TeatimeMessage{
		Recipient: this.PeerList[peerID].RepoRemoteName,
		Action:    encoder.ACTION_PING,
		Payload: encoder.PingPayload{
			PingID:         this.NumPingsSent,
			CurrentRetries: 0, // Value is not used
			IsPong:         false,
			OriginIP:       "", // Set by receiving Crumpet
		},
	}
	serializer := encoder.InterTeatimeSerializer{}
	encoded, err := serializer.Serialize(msg)
	if err != nil {
		return err
	}

	peerConnection, err := makeTCPConn(fmt.Sprintf("%v:%v", peerID, tt.TEATIME_DEFAULT_PORT))
	if err != nil {
		return err
	}
	defer peerConnection.Close()

	this.NumPingsSent++
	_, err = tt.SendData(peerConnection, encoded)
	if err != nil {
		return err
	}

	return nil
}

func (this *TTNetSession) sendTTPong(peerID string) error {
	msg := encoder.TeatimeMessage{
		Recipient: this.PeerList[peerID].RepoRemoteName,
		Action:    encoder.ACTION_PING,
		Payload: encoder.PingPayload{
			PingID:         this.NumPongsSent,
			CurrentRetries: 0, // Value is not used
			IsPong:         true,
			OriginIP:       "", // Set by receiving Crumpet
		},
	}
	serializer := encoder.InterTeatimeSerializer{}
	encoded, err := serializer.Serialize(msg)
	if err != nil {
		return err
	}

	peerConnection, err := makeTCPConn(fmt.Sprintf("%v:%v", peerID, tt.TEATIME_DEFAULT_PORT))
	if err != nil {
		return err
	}
	defer peerConnection.Close()

	this.NumPongsSent++
	_, err = tt.SendData(peerConnection, encoded)
	if err != nil {
		return err
	}

	return nil
}

func (this *TTNetSession) startPingService(peerID string, pingInterval time.Duration, repoRemoteName string) {
	ticker := time.NewTicker(pingInterval)

	for {
		log.Printf("[Repo: %v] Pinging %v at %v (id=%d)", this.Repo.Name, peerID, time.Now(), this.NumPingsSent)
		this.sendTTPing(peerID)

		<-ticker.C

		if this.NumPingsSent > 20 && this.NumPingsSent > 2*this.NumPongsRcvd {
			log.Printf("Unreliable connection. Closing ", peerID)

			// Remove peer
			delete(this.PeerList, peerID)

			break
		}
	}

	// If this is ever reached... we should attempt a reconnection.
	// NOTE: This is NOT just a TCP reconnect. TryTeaTimeConnect sends a CONNECT
	// request payload. Assuming the peer is still up, they may have removed us
	// as a peer, and are therefore rejecting our pings. Hence, we should resend
	// connect request.
	if err := this.TryTeaTimeConn(peerID, repoRemoteName); err != nil {
		log.Printf("Reconnect failed. Aborting.")
	}
}
