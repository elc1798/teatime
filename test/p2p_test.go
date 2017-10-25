package test

import (
    "testing"
    "time"

    p2p "github.com/elc1798/teatime/p2p"
)

func TestBasicServer(t *testing.T) {
    // Set up
    p2p.StartListener(12345)
    testSession := p2p.NewTTNetSession()

    timer := time.NewTimer(time.Millisecond * 250)

    err := testSession.TryTeaTimeConn("localhost:12345")

    <- timer.C

    if err != nil {
        t.Fatalf("Error in TryTeaTimeConn: %v\n", err)
    }

    if len(testSession.PeerConns) != 1 {
        t.Fatalf("Failed to append to peer connections: %v\n", testSession.PeerConns)
    }

    if len(testSession.PeerList) != 1 {
        t.Fatalf("Failed to append to peer list: %v\n", testSession.PeerList)
    }
}
