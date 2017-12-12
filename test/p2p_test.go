package test

import (
	"os"
	"reflect"
	"testing"
	"time"

	tt "github.com/elc1798/teatime"
	crumpet "github.com/elc1798/teatime/crumpet"
	p2p "github.com/elc1798/teatime/p2p"
)

const REPO_1 = "tt_test1"
const REPO_2 = "tt_test2"

func TestBasicServer(t *testing.T) {
	tt.ResetTeatime()
	p2p.PING_INTERVAL = time.Millisecond * 250

	r1, d1, _ := setUpRepos(REPO_1)
	defer os.RemoveAll(d1)

	r2, d2, _ := setUpRepos(REPO_2)
	defer os.RemoveAll(d2)

	daemon := crumpet.NewCrumpetDaemon()
	daemon.Start(false)

	session1 := p2p.NewTTNetSession(r1)
	session2 := p2p.NewTTNetSession(r2)

	// Connect session1 to session2
	timer := time.NewTimer(time.Millisecond * 500)
	start_time := time.Now()
	err := session1.TryTeaTimeConn("localhost", r2.Name)
	t.Logf("TryTeaTimeConn took %v", time.Since(start_time))
	<-timer.C
	if err != nil {
		t.Fatalf("TryTeaTimeConn failed: %v", err)
	}

	if len(session1.PeerList) != 1 {
		t.Fatalf("Failed to append to peer list: %v\n", session1.PeerList)
	}

	if len(session2.PeerList) != 1 {
		t.Fatalf("Failed to append to peer list: %v\n", session2.PeerList)
	}

	if session1.PeerList["127.0.0.1"].IP != "127.0.0.1" {
		t.Fatalf("Invalid peer connection: IP=%v\n", session1.PeerList["127.0.0.1"].IP)
	}

	if session1.PeerList["127.0.0.1"].Port != tt.TEATIME_DEFAULT_PORT {
		t.Fatalf("Invalid peer connection: Port=%v\n", session1.PeerList["127.0.0.1"].Port)
	}

	if session1.PeerList["127.0.0.1"].RepoRemoteName != r2.Name {
		t.Fatalf("Invalid peer: %v", session1.PeerList["127.0.0.1"])
	}

	if session2.PeerList["127.0.0.1"].RepoRemoteName != r1.Name {
		t.Fatalf("Invalid peer: %v", session2.PeerList["127.0.0.1"])
	}

	// Check that there were 2 ping pongs each
	expected := map[string][2]int{
		"NumPingsSent": [2]int{2, 2},
		"NumPingsRcvd": [2]int{2, 2},
		"NumPongsSent": [2]int{2, 2},
		"NumPongsRcvd": [2]int{2, 2},
	}

	getValue := func(sess *p2p.TTNetSession, field string) int {
		r1 := reflect.ValueOf(sess)
		f1 := reflect.Indirect(r1).FieldByName(field)
		return int(f1.Int())
	}

	for field, sol := range expected {
		v0 := getValue(session1, field)
		v1 := getValue(session2, field)

		t.Logf("Session1 has %d of %d %s", v0, sol[0], field)
		t.Logf("Session2 has %d of %d %s", v1, sol[1], field)

		if v0 < sol[0] {
			t.Fatalf("Session1 has %d of %d %s", v0, sol[0], field)
		}

		if v1 < sol[1] {
			t.Fatalf("Session2 has %d of %d %s", v1, sol[1], field)
		}
	}

	os.RemoveAll(tt.TEATIME_DEFAULT_HOME)
}
