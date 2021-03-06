package test

import (
	"os"
	"testing"

	tt "github.com/elc1798/teatime"
	p2p "github.com/elc1798/teatime/p2p"
)

const REPO = "tt_test"

func TestPeerCache(t *testing.T) {
	// Clear original teatime directory
	tt.ResetTeatime()

	r1, d1, _ := setUpRepos(REPO)

	testPeerList := map[string]p2p.Peer{
		"1.2.3.4:80":         p2p.Peer{IP: "1.2.3.4", Port: 80, RepoRemoteName: "a"},
		"8.8.8.8:443":        p2p.Peer{IP: "8.8.8.8", Port: 443, RepoRemoteName: "bob"},
		"172.217.9.46:12345": p2p.Peer{IP: "172.217.9.46", Port: 12345, RepoRemoteName: "joe"},
		"192.168.1.1:9001":   p2p.Peer{IP: "192.168.1.1", Port: 9001, RepoRemoteName: "lmao"},
		"192.17.13.36:1111":  p2p.Peer{IP: "192.17.13.36", Port: 1111, RepoRemoteName: "xd"},
	}

	testSession := p2p.NewTestTTNetSession(r1)
	testSession.PeerList = testPeerList
	if err := testSession.GenerateLocalPeerCache(); err != nil {
		t.Fatalf("Error creating peer cache: %v\n", err)
	}

	readPeerList, err := testSession.GetLocalPeerCache()
	if err != nil {
		t.Fatalf("Error reading local cache: %v\n", err)
	}

	if len(testPeerList) != len(readPeerList) {
		t.Fatalf("Invalid peer list!")
	}

	for i, peer := range testPeerList {
		if peer.IP != readPeerList[i].IP || peer.Port != readPeerList[i].Port || peer.RepoRemoteName != readPeerList[i].RepoRemoteName {
			t.Fatalf("Invalid peer list!")
		}
	}

	os.RemoveAll(tt.TEATIME_DEFAULT_HOME)
	os.RemoveAll(d1)
}
