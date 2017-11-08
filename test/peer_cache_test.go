package test

import (
	"os"
	"testing"

	tt "github.com/elc1798/teatime"
	p2p "github.com/elc1798/teatime/p2p"
)

func TestPeerCache(t *testing.T) {
	os.RemoveAll(tt.TEATIME_DEFAULT_HOME)
	if err := os.Mkdir(tt.TEATIME_DEFAULT_HOME, 0777); err != nil {
		t.Fatalf("Error creating teatime home: %v\n", err)
	}

	testPeerList := map[string]p2p.Peer{
		"1.2.3.4:80":         p2p.Peer{IP: "1.2.3.4", Port: 80},
		"8.8.8.8:443":        p2p.Peer{IP: "8.8.8.8", Port: 443},
		"172.217.9.46:12345": p2p.Peer{IP: "172.217.9.46", Port: 12345},
		"192.168.1.1:9001":   p2p.Peer{IP: "192.168.1.1", Port: 9001},
		"192.17.13.36:1111":  p2p.Peer{IP: "192.17.13.36", Port: 1111},
	}

	if err := p2p.GenerateLocalPeerCache(testPeerList); err != nil {
		t.Fatalf("Error creating peer cache: %v\n", err)
	}

	readPeerList, err := p2p.GetLocalPeerCache()
	if err != nil {
		t.Fatalf("Error reading local cache: %v\n", err)
	}

	if len(testPeerList) != len(readPeerList) {
		t.Fatalf("Invalid peer list!")
	}

	for i, peer := range testPeerList {
		if peer.IP != readPeerList[i].IP || peer.Port != readPeerList[i].Port {
			t.Fatalf("Invalid peer list!")
		}
	}

	os.RemoveAll(tt.TEATIME_DEFAULT_HOME)
}
