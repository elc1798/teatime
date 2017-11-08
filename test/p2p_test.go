package test

import (
	"reflect"
	"testing"
	"time"

	p2p "github.com/elc1798/teatime/p2p"
)

func TestBasicServer(t *testing.T) {
	// Set up
	pingInterval := time.Millisecond * 150

	serverSession := p2p.NewTTNetSession()
	serverSession.StartListener(12345, false)

	testSession := p2p.NewTTNetSession()
	timer := time.NewTimer(time.Millisecond * 300)

	start_time := time.Now()
	err := testSession.TryTeaTimeConn("localhost:12345", pingInterval)
	t.Logf("Connection took %v", time.Since(start_time))

	<-timer.C

	if err != nil {
		t.Fatalf("Error in TryTeaTimeConn: %v\n", err)
	}

	if len(testSession.PeerConns) != 1 {
		t.Fatalf("Failed to append to peer connections: %v\n", testSession.PeerConns)
	}

	if len(testSession.PeerList) != 1 {
		t.Fatalf("Failed to append to peer list: %v\n", testSession.PeerList)
	}

	if len(serverSession.PeerList) != 1 {
		t.Fatalf("Failed to append to peer list: %v\n", serverSession.PeerList)
	}

	// Check for peers. The server session should have testSession as peer, and
	// vice versa
	if testSession.PeerList["127.0.0.1:12345"].IP != "127.0.0.1" {
		t.Fatalf("Invalid peer connection: IP=%v\n", testSession.PeerList["127.0.0.1:12345"].IP)
	}

	if testSession.PeerList["127.0.0.1:12345"].Port != 12345 {
		t.Fatalf("Invalid peer connection: Port=%v\n", testSession.PeerList["127.0.0.1:12345"].Port)
	}

	// Test Session Local Addr should be equal to server remote
	for _, v := range serverSession.PeerList {
		if v.IP != "127.0.0.1" {
			t.Fatalf("Invalid peer connection: %v\n", v)
		}
	}

	// Check that there were 2 ping pongs each
	expected := map[string][2]int{
		"NumPingsSent": [2]int{2, 0},
		"NumPingsRcvd": [2]int{0, 2},
		"NumPongsSent": [2]int{0, 2},
		"NumPongsRcvd": [2]int{2, 0},
	}

	getValue := func(sess *p2p.TTNetSession, field string) int {
		r1 := reflect.ValueOf(sess)
		f1 := reflect.Indirect(r1).FieldByName(field)
		return int(f1.Int())
	}

	for field, sol := range expected {
		v0 := getValue(testSession, field)
		v1 := getValue(serverSession, field)

		t.Logf("TestSession has %d of %d %s", v0, sol[0], field)
		t.Logf("ServerSession has %d of %d %s", v1, sol[1], field)

		if v0 != sol[0] {
			t.Fatalf("TestSession has %d of %d %s", v0, sol[0], field)
		}

		if v1 != sol[1] {
			t.Fatalf("ServerSession has %d of %d %s", v1, sol[1], field)
		}
	}
}
